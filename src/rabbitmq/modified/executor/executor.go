package executor

type Executor struct{}

func NewExecutor() Executor {
	return Executor{}
}

/*
func Execute (e *Executor){

case <-s.ChStart:
count = 0
t1 = time.Now()
case <-s.ChStop:
// calculate arrival rate
monitorInterval := time.Now().Sub(t1).Seconds()
numberOfMessages := float64(count)

s.ArrivalRate = numberOfMessages / monitorInterval

//runtime.ReadMemStats(&m) // TODO

// inspect queue
q1, err1 := s.Ch.QueueInspect("rpc_queue")
shared.FailOnError(err1, "Failed to inspect the queue")

//if s.ArrivalRate == 0.0 { // finished processing
//	os.Exit(0)
//} else {
//fmt.Printf("%d;%d;%.3f \n", u, q1.Messages, s.ArrivalRate)
//}

fmt.Printf("%d;%d;%.3f \n", s.PC, q1.Messages, s.ArrivalRate)

if !s.IsAdaptive { // for experimental purpose
if countSample < 29 {
countSample++
} else {
//s.PC = s.PC + 1000
countSample = 0

// set qos
err := s.Ch.Qos(
s.PC, // update prefetch count
0,    // prefetch size
true, // default is false
)
shared.FailOnError(err, "Failed to set QoS")
}
}

// Reconfigure QoS (Ineffective if autoAck is true)
if s.IsAdaptive {

// compute new value of control law
u = int(controller.Update(s.Ctler, s.ArrivalRate))

// queue controller
//u = int(controller.Update(s.Ctler, float64(q1.Messages)))

// set qos
err := s.Ch.Qos(
u,    // update prefetch count
0,    // prefetch size
true, // default is false
)
shared.FailOnError(err, "Failed to set QoS")
}
}

*/
