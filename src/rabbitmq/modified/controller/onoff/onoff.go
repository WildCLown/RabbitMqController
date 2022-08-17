package onoff

//import "fmt"

type OnOff struct {
	R     float64
	U     float64
	Y     float64
	Error float64
}

func (c *OnOff) InitC(p ...float64) {
	r := p[0]

	c.R = r
	c.U = 0.0
	c.Y = 0.0
	c.Error = 0.0
}

func (c *OnOff) Update(p ...float64) float64 {
	y := p[0]
	c.Y = y

	c.Error = c.R - y

	if int(c.Error) >= 0 { // error = r - y
		c.U = 1000
	}

	if c.Error < 0 {
		c.U = 1
	}

	return c.U
}
