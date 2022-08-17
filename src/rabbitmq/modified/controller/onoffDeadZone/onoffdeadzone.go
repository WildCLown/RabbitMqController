package onoffDeadZone

//import "fmt"

type OnOffDeadZone struct {
	R     float64
	U     float64
	Y     float64
	Error float64
}

func (c *OnOffDeadZone) InitC(p ...float64) {
	r := p[0]

	c.R = r
	c.U = 1.0 // fix this problem by properly setting the prefetch
	c.Y = 0.0
	c.Error = 0.0
}

func (c *OnOffDeadZone) Update(p ...float64) float64 {
	y := p[0]
	c.Y = y

	c.Error = c.R - y
	percentage := c.Error / c.R * 100.0

	//fmt.Println(percentage)

	if percentage <= -20 { // error = r - y consider 10%
		c.U = 1
	}

	if (percentage) >= 20 {
		c.U = 100
	}

	return c.U
}
