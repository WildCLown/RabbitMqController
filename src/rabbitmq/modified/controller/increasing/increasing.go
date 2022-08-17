package increasing

type Increasing struct {
	R     float64
	U     float64
	Y     float64
	Error float64
}

func (c *Increasing) InitC(p ...float64) {
	r := p[0]

	c.R = r
	c.U = 0.0
	c.Y = 0.0
	c.Error = 0.0
}

func (c *Increasing) Update(p ...float64) float64 {
	c.U = c.U+25

	return c.U
}
