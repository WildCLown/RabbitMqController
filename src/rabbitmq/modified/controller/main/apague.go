package main

import (
	"fmt"
	"math/rand"
	"rabbitmq/modified/controller/pid"
)

func main() {
    c := pid.PIDController{}

	c.Init(1.0,0.0,0.0)

	sp := 1600.0

	for i := 0; i < 100; i++{
		m := float64(rand.Intn(2000))
		fmt.Printf("%.2f;%.2f\n",m,c.Update(sp,m))
	}
}
