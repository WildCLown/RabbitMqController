package pid

import (
	"rabbitmq/acmsac/controller/info"
)

type PIDController struct {
	Info info.InfoController
}

func (c *PIDController) Initialise(p ...float64) {

	r := p[0] // reference input

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

	c.Info.Integrator = 0.0

	c.Info.PreviousError = 0.0
	c.Info.SumPreviousErrors = 0.0
	c.Info.Out = 0.0
}

func (c *PIDController) Update(p ...float64) float64 {
	//var limMinInt, limMaxInt float64

	y := p[0]

	// errors
	err := c.Info.R - y

	// Proportional (David page 49)
	proportional := c.Info.Kp * err // pcontroller[k] = kp * e[k]

	// Integrator (David page 49)
	c.Info.Integrator += 1.0 * err // TODO
	integrator := c.Info.Integrator * c.Info.Ki

	// Differentiator (David page 49)
	differentiator := c.Info.Kd * (err - c.Info.PreviousError) / 1.0 // TODO

	// pid output
	c.Info.Out = proportional + integrator + differentiator

	if c.Info.Out > c.Info.LimMax {
		c.Info.Out = c.Info.LimMax
	} else if c.Info.Out < c.Info.LimMin {
		c.Info.Out = c.Info.LimMin
	}

	c.Info.PreviousError = err
	c.Info.SumPreviousErrors = c.Info.SumPreviousErrors + err

	return c.Info.Out
}
