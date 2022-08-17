package pi

type PIController struct {
	R  float64
	Kp float64
	Ki float64
	Kd float64

	Tau float64 // 1  TODO

	LimMin float64 // 1
	LimMax float64 // 10000

	T float64 // 1 TODO

	Integrator          float64
	Differentiator      float64
	PreviousError       float64
	SumPreviousErrors   float64
	PreviousMeasurement float64
	Out                 float64
}

func (c *PIController) Initialise(p ...float64) {

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
	c.PreviousMeasurement = 0.0
	c.Out = 0.0

	c.T = 1.0   // TODO
	c.Tau = 1.0 // TODO rabbitmq

	//pid.T = 0.1     // test
	//pid.Tau = 0.9 // test
}

func (c *PIController) Update(p ...float64) float64 {
	//var limMinInt, limMaxInt float64

	y := p[0]

	// error
	error := c.R - y

	// Proportional
	//proportional := pid.Kp * error // pcontroller[k] = kp * e[k]

	// Integrator i[k] = i[k-1] + ()
	//pid.Integrator = pid.Integrator + 0.5*pid.Ki*pid.T*(error+pid.PreviousError)

	//u(k) = u(k − 1) + (KP + KI)e(k) − KPe(k − 1) Page 302
	c.Out = c.Out + (c.Kp+c.Ki)*error - c.Kp*c.PreviousError

	// Et =δt ·et +Et−1  ui,t =kiEt page 44 (David) -- integral part
	//c.Out = c.Kp*error + c.Ki*(1.0*error+c.SumPreviousErrors) // TODO 10 is 10 s

	if c.Out > c.LimMax {
		c.Out = c.LimMax
	} else if c.Out < c.LimMin {
		c.Out = c.LimMin
	}

	c.SumPreviousErrors = c.SumPreviousErrors + error
	c.PreviousError = error
	c.PreviousMeasurement = y

	//fmt.Println("****** Controller::", c.R, c.Kp, c.Ki, c.Kd, c.Out)

	return c.Out
}
