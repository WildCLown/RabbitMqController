package pcontroller

type PController struct {
	R  float64
	Kp float64

	LimMin float64 // 1
	LimMax float64 // 10000

	Out float64
}

func (c *PController) InitC(p ...float64) {

	r := p[0] // reference input

	kp := p[1]
	//ki := p[2] not used
	//kd := p[3] not used
	limMin := p[4]
	limMax := p[5]

	c.R = r
	c.LimMin = limMin
	c.LimMax = limMax

	c.Kp = kp

	c.Out = 0.0
}

func (c *PController) Update(p ...float64) float64 {
	//var limMinInt, limMaxInt float64

	y := p[0]

	// errors
	error := c.R - y

	// Proportional
	proportional := c.Kp * error // pcontroller[k] = kp * e[k]

	// pcontroller output
	c.Out = proportional

	if c.Out > c.LimMax {
		c.Out = c.LimMax
	} else if c.Out < c.LimMin {
		c.Out = c.LimMin
	}

	//fmt.Println(c.LimMin, c.LimMax, c.Kp, c.Out)
	return c.Out
}
