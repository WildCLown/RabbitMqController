package main

import (
	"fmt"
	"math/rand"
	"rabbitmq/modified/controller/pid"
)

func main() {

	c := pid.PIDController{}
	r := 5000.0
	kp := 0.0056433669
	ki := -0.0204420019
	kd := 0.001412706
	min := 1.0
	max := 10000.0

	c.InitC(r, kp, ki, kd, min, max)

	for i := 0; i < 100; i++ {
		y := float64(rand.Intn(10000))
		u := c.Update(y)
		fmt.Println(y, u)
	}
}
