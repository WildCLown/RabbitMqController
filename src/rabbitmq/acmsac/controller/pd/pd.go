package pd

import (
	"rabbitmq/acmsac/controller/info"
)

type PDController struct {
	Info info.InfoController
}

func (c *PDController) Initialise(p ...float64) {

	r := p[0]

	kp := p[1]
	ki := p[2]
	kd := p[3]

	limMin := p[4]
	limMax := p[5]

	c.Info.R = r
	c.Info.LimMin = limMin
	c.Info.LimMax = limMax

	c.Info.Kp = kp
	c.Info.Ki = ki
	c.Info.Kd = kd
}

func (c *PDController) Update(p ...float64) float64 {
	//var limMinInt, limMaxInt float64

	y := p[0]

	// errors
	err := c.Info.R - y

	//u(k) = KPe(k) + KD(e(k) − e(k − 1)) // page 316
	c.Info.Out = c.Info.Kp*err + c.Info.Kd*(err-c.Info.PreviousError)

	// pd output page 46 (David)
	//ud,t =kd (et −et−1)/Tetat, where Tetat is in time unit
	//c.Info.Out = c.Info.Kd * (error - c.Info.PreviousError) / 10.0 // TODO

	// pd output
	//c.Out = proportional + c.Integrator + c.Differentiator // page 316

	if c.Info.Out > c.Info.LimMax {
		c.Info.Out = c.Info.LimMax
	} else if c.Info.Out < c.Info.LimMin {
		c.Info.Out = c.Info.LimMin
	}

	c.Info.PreviousError = err
	c.Info.SumPreviousErrors = c.Info.SumPreviousErrors + err

	//fmt.Println("pd::", y, c.R, c.Kp, c.Ki, c.Kd, c.SumPreviousErrors, c.Out)

	return c.Info.Out
}
