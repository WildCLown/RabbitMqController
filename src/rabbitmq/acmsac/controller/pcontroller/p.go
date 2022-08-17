package pcontroller

import (
	"rabbitmq/acmsac/controller/info"
)

type PController struct {
	Info info.InfoController
}

func (c *PController) Initialise(p ...float64) {

	r := p[0]

	kp := p[1]
	//ki := p[2] not used
	//kd := p[3] not used
	limMin := p[4]
	limMax := p[5]

	c.Info.R = r
	c.Info.LimMin = limMin
	c.Info.LimMax = limMax

	c.Info.Kp = kp

	c.Info.Out = 0.0
}

func (c *PController) Update(p ...float64) float64 {
	y := p[0]

	// errors
	error := c.Info.R - y

	// Proportional
	proportional := c.Info.Kp * error // pcontroller[k] = kp * e[k]

	c.Info.Out = proportional

	if c.Info.Out > c.Info.LimMax {
		c.Info.Out = c.Info.LimMax
	} else if c.Info.Out < c.Info.LimMin {
		c.Info.Out = c.Info.LimMin
	}

	return c.Info.Out
}
