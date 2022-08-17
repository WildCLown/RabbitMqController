package pid

type PIDController struct {
	R  float64
	Kp float64
	Ki float64
	Kd float64

	LimMin float64 // 1
	LimMax float64 // 10000

	Integrator        float64
	Differentiator    float64
	SumPreviousErrors float64
	PreviousError     float64
	Out               float64
}

func (c *PIDController) InitC(p ...float64) {

	r := p[0] // reference input

	kp := p[1]
	ki := p[2]
	kd := p[3]

	limMin := p[4]
	limMax := p[5]

	c.R = r
	c.LimMin = limMin
	c.LimMax = limMax

	c.Kp = kp
	c.Ki = ki
	c.Kd = kd

	c.Integrator = 0.0
	c.PreviousError = 0.0
	c.SumPreviousErrors = 0.0
	c.Differentiator = 0.0
	c.Out = 0.0
}

func (c *PIDController) Update(p ...float64) float64 {
	//var limMinInt, limMaxInt float64

	y := p[0]

	// errors
	error := c.R - y

	// Proportional
	//proportional := c.Kp * error // pcontroller[k] = kp * e[k]

	// Integrator i[k] = Ki * Sum[i=0,k-1](ei) + Ki*e(k)
	//c.Integrator = c.Ki*c.SumPreviousErrors + c.Ki*error

	//if c.LimMax > proportional {
	//	limMaxInt = c.LimMax - proportional
	//} else {
	//	limMaxInt = 0.0
	//}

	//if c.LimMin < proportional {
	//	limMinInt = c.LimMin - proportional
	//} else {
	//	limMinInt = 0.0
	//}

	//if c.Integrator > limMaxInt {
	//	c.Integrator = limMaxInt
	//} else if c.Integrator < limMinInt {
	//	c.Integrator = limMinInt
	//}

	// Differentiator
	//c.Differentiator = c.Kd * (error - c.PreviousError)

	//u(k) = KPe(k) + KISum(k−1, i=0)e(i) + KIe(k) + KD [e(k) − e(k − 1)] // page 320
	c.Out = c.Kp*error + c.Ki*c.SumPreviousErrors + c.Ki*error + c.Kd*(error-c.PreviousError)

	// See page 48 (David)
	//c.Out = c.Kp*error + c.Ki*(10.0*error+c.SumPreviousErrors) + c.Kd*(error-c.PreviousError) //TODO 10.0

	//fmt.Println("PID01", c.Out, c.LimMin, c.LimMax, c.SumPreviousErrors)

	// pid output
	//c.Out = proportional + c.Integrator + c.Differentiator // page 320

	if c.Out > c.LimMax {
		c.Out = c.LimMax
	} else if c.Out < c.LimMin {
		c.Out = c.LimMin
	}

	//if c.Out > c.LimMax {
	//	c.Out = c.LimMax
	//} else if c.Out < c.LimMin {
	//	c.Out = c.LimMin
	//}

	c.PreviousError = error
	c.SumPreviousErrors = c.SumPreviousErrors + error

	//fmt.Println("PID::", y, c.R, c.Kp, c.Ki, c.Kd, c.SumPreviousErrors, c.Out)

	return c.Out
}
