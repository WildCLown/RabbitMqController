package onoff

import (
	"rabbitmq/acmsac/controller/info"
)

type OnOff struct {
	Info info.InfoController
}

func (c *OnOff) Initialise(p ...float64) {
	c.Info.R = p[0]
	c.Info.LimMin = p[1]
	c.Info.LimMax = p[2]
}

func (c *OnOff) Update(p ...float64) float64 {
	y := p[0]

	error := c.Info.R - y

	if int(error) > 0 {
		c.Info.Out = c.Info.LimMax
	}

	if error < 0 {
		c.Info.Out = c.Info.LimMin
	}

	return c.Info.Out
}
