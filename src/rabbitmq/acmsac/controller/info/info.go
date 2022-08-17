package info

type InfoController struct {
	R  float64
	Kp float64
	Ki float64
	Kd float64

	LimMin float64
	LimMax float64

	Integrator        float64
	SumPreviousErrors float64
	PreviousError     float64
	Out               float64
}
