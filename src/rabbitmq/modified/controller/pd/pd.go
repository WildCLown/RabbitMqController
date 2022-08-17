package pd

type PDController struct {
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

func (c *PDController) InitC(p ...float64) {

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

func (c *PDController) Update(p ...float64) float64 {
	//var limMinInt, limMaxInt float64

	y := p[0]

	// errors
	error := c.R - y

	//u(k) = KPe(k) + KD(e(k) − e(k − 1)) // page 316
	c.Out = c.Kp*error + c.Kd*(error-c.PreviousError)

	// pd output page 46 (David)
	//ud,t =kd (et −et−1)/Tetat, where Tetat is in time unit
	c.Out = c.Kd * (error - c.PreviousError) / 10.0 // TODO

	// pd output
	//c.Out = proportional + c.Integrator + c.Differentiator // page 316

	if c.Out > c.LimMax {
		c.Out = c.LimMax
	} else if c.Out < c.LimMin {
		c.Out = c.LimMin
	}

	c.PreviousError = error
	c.SumPreviousErrors = c.SumPreviousErrors + error

	//fmt.Println("pd::", y, c.R, c.Kp, c.Ki, c.Kd, c.SumPreviousErrors, c.Out)

	return c.Out
}
