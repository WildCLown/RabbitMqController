package controller

import (
	"fmt"
	"os"
	"rabbitmq/acmsac/controller/onoff"
	"rabbitmq/acmsac/controller/onoffDeadZone"
	"rabbitmq/acmsac/controller/pid"
	"rabbitmq/shared"
)

type IController interface {
	Update(p ...float64) float64
	Initialise(p ...float64)
}

func NewController(typeName string, p ...float64) IController {
	switch typeName {
	case "OnOff":
		r := p[0]
		limMin := shared.PC_DEFAULT_LIMIT_MIN
		limMax := shared.PC_DEFAULT_LIMIT_MAX
		c := onoff.OnOff{}
		c.Initialise(r, float64(limMin), float64(limMax))
		return &c
	case "OnOffDeadZone":
		r := p[0]
		limMin := shared.PC_DEFAULT_LIMIT_MIN
		limMax := shared.PC_DEFAULT_LIMIT_MAX
		c := onoffDeadZone.OnOffDeadZone{}
		c.Initialise(r, float64(limMin), float64(limMax))
		return &c
	case "PID":
		r := p[0]
		kp := p[1]
		ki := p[2]
		kd := p[3]
		limMin := shared.PC_DEFAULT_LIMIT_MIN
		limMax := shared.PC_DEFAULT_LIMIT_MAX
		c := pid.PIDController{}
		c.Initialise(r, kp, ki, kd, float64(limMin), float64(limMax))
		return &c
	default:
		fmt.Println("Error: Controller type unknown")
		os.Exit(0)
	}

	return *new(IController)
}

func Update(c IController, y float64) float64 {

	return c.Update(y)
}

func Initialise(c IController) {
	c.Initialise()
}
