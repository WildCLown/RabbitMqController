package controller

import (
	"fmt"
	"os"
	"rabbitmq/modified/controller/onoff"
	"rabbitmq/modified/controller/onoffDeadZone"
	"rabbitmq/modified/controller/pcontroller"
	"rabbitmq/modified/controller/pd"
	"rabbitmq/modified/controller/pi"
	"rabbitmq/modified/controller/pid"
)

type IController interface {
	Update(p ...float64) float64
	InitC(p ...float64)
}
type Controller struct {
	TypeName string
	R        float64 // setpoint
	Type     interface{}
}

func NewController(typeName string, p ...float64) IController {
	x := new(IController)

	switch typeName {
	case "OnOffDeadZone":
		r := p[0]
		//kp := p[1] // not used
		//ki := p[2]
		//kd := p[3]
		//limMin := p[4]
		//limMax := p[5]
		c := onoffDeadZone.OnOffDeadZone{}
		c.InitC(r)
		return &c
	case "OnOff":
		r := p[0]
		//kp := p[1] // not used
		//ki := p[2]
		//kd := p[3]
		//limMin := p[4]
		//limMax := p[5]
		c := onoff.OnOff{} // TODO
		c.InitC(r)
		return &c
	case "P":
		r := p[0]
		kp := p[1]
		ki := p[2]
		kd := p[3]
		limMin := p[4]
		limMax := p[5]
		c := pcontroller.PController{} // TODO
		c.InitC(r, kp, ki, kd, limMin, limMax)
		return &c
	case "PI":
		r := p[0]
		kp := p[1]
		ki := p[2]
		kd := p[3]
		limMin := p[4]
		limMax := p[5]
		c := pi.PIController{} // TODO
		c.InitC(r, kp, ki, kd, limMin, limMax)
		return &c
	case "PD":
		r := p[0]
		kp := p[1]
		ki := p[2]
		kd := p[3]
		limMin := p[4]
		limMax := p[5]
		c := pd.PDController{} // TODO
		c.InitC(r, kp, ki, kd, limMin, limMax)
		return &c
	case "PID":
		r := p[0]
		kp := p[1]
		ki := p[2]
		kd := p[3]
		limMin := p[4]
		limMax := p[5]
		c := pid.PIDController{} // TODO
		c.InitC(r, kp, ki, kd, limMin, limMax)
		return &c
	default:
		fmt.Println("Error: Controller type unknown")
		os.Exit(0)
	}

	return *x
}

func Update(c IController, y float64) float64 {

	return c.Update(y)
}

func InitC(c IController) {
	c.InitC()
}
