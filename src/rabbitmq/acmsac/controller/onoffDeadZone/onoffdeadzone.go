package onoffDeadZone

import (
	"rabbitmq/acmsac/controller/info"
)

//import "fmt"

type OnOffDeadZone struct {
	Info info.InfoController
}

func (c *OnOffDeadZone) Initialise(p ...float64) {
	c.Info.R = p[0]
	c.Info.LimMin = p[1]
	c.Info.LimMax = p[2]
}

func (c *OnOffDeadZone) Update(p ...float64) float64 {
	y := p[0]

	// fix deadzone to 20%
	upper := c.Info.R * 0.20
	lower := -1 * c.Info.R * 0.20

	error := c.Info.R - y

	if error > 0 && error > upper {
		c.Info.Out = c.Info.LimMax
	}

	if error < 0 && error < lower {
		c.Info.Out = c.Info.LimMax
	}

	return c.Info.Out
}
