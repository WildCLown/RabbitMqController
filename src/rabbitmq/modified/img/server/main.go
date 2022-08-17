package main

import (
	"flag"
	"fmt"
	"github.com/streadway/amqp"
	_ "net/http/pprof"
	"os"
	"rabbitmq/modified/controller"
	"rabbitmq/modified/executor"
	"rabbitmq/modified/monitor"
	"rabbitmq/shared"
	"runtime"
	"time"
)

type Server struct {
	IsAdaptive      bool
	MonitorInterval time.Duration
	Conn            *amqp.Connection
	Ch              *amqp.Channel
	Queue           amqp.Queue
	Msgs            <-chan amqp.Delivery
	ChStart         chan bool
	ChStop          chan bool
	Mnt             monitor.Monitor
	Ctler           controller.IController
	ArrivalRate     float64
	PC              int
	Exe             executor.Executor
}

func main() {

	// configure/read flags
	var isAdaptivePtr = flag.Bool("is-adaptive", false, "is-adaptive is a boolean")
	var controllerTypePtr = flag.String("controller-type", "OnOff", "controller-type is a string")
	var monitorIntervalPtr = flag.Int("monitor-interval", 1, "monitor-interval is an int (s)")
	var setPointPtr = flag.Float64("set-point", 3000.0, "set-point is a float (goal rate)")
	var kpPtr = flag.Float64("kp", 1.0, "Kp is a float")
	var kiPtr = flag.Float64("ki", 1.0, "Ki is a float")
	var kdPtr = flag.Float64("kd", 1.0, "Kd is a float")
	var prefetchCountPtr = flag.Int("prefetch-count", 1, "prefetch-count is an int")
	flag.Parse()

	var c controller.IController
	r := *setPointPtr

	// create controller
	c = controller.NewController(*controllerTypePtr, r, *kpPtr, *kiPtr, *kdPtr, 1.0, 10000.0)

	// create new consumer
	var server = NewServer(*isAdaptivePtr, *monitorIntervalPtr, c, *prefetchCountPtr)

	fmt.Println("Server started [", *isAdaptivePtr, *controllerTypePtr, "Kp=", *kpPtr, "Ki=", *kiPtr, "Kd=", *kdPtr, "Goal=", r, "Monitor Interval=", *monitorIntervalPtr, "PC=", *prefetchCountPtr, "]")

	// run consumer
	server.Run() // todo
}

func NewServer(isAdaptive bool, monitorInterval int, c controller.IController, prefetchCount int) Server {
	s := Server{}

	// Configure consumer
	s.IsAdaptive = isAdaptive
	s.MonitorInterval = time.Duration(monitorInterval) * time.Second

	// Initialise channel to communicate with Monitor
	s.ChStart = make(chan bool)
	s.ChStop = make(chan bool)

	// create Monitor
	s.Mnt = monitor.NewMonitor(s.MonitorInterval)

	// create analyser

	//create Planner

	// create Executor
	s.Exe = executor.NewExecutor()

	// set controller
	s.Ctler = c

	// set initial PC -- always 1
	s.PC = prefetchCount

	return s
}

// Run consumer
func (s Server) Run() {

	// close all rabbitmq elements before exiting
	defer s.Conn.Close()
	defer s.Ch.Close()

	// start monitor
	go s.Mnt.Monitoring(s.ChStart, s.ChStop)

	// Configure RabbitMQ
	s.configureRabbitMQ()

	// handle requests
	s.handleRequests()
}

// Handle requests
func (s Server) handleRequests() {

	//s.PC = 145 // TODO experimental purpose
	u := s.PC // initial u
	count := 0
	countSample := 0

	t1 := time.Time{}
	//for {
	for d := range s.Msgs {
		//select {
		//case d := <-s.Msgs: // receive a message
		// configure App business time
		time.Sleep(1 * time.Millisecond) // TODO for experimental purposes
		d.Ack(false)
		count++
		select {
		case <-s.ChStart: // monitor
			count = 0
			t1 = time.Now()
		case <-s.ChStop: // monitor
			// calculate arrival rate
			monitorInterval := time.Now().Sub(t1).Seconds()
			numberOfMessages := float64(count)

			//fmt.Println("Monitor Interval=", monitorInterval, "Number of Messages=", numberOfMessages)

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

			// Non-adaptive
			if !s.IsAdaptive { // for experimental purpose
				if countSample < 99 {
					countSample++
				} else {
					s.PC = s.PC + 1
					countSample = 0

					// set qos
					err := s.Ch.Qos(
						s.PC, // update prefetch count
						0,    // prefetch size
						true, // default is false
					)
					shared.FailOnError(err, "Failed to set QoS")
				}
			} else { // adaptive
				// Reconfigure QoS (Ineffective if autoAck is true)
				// compute new value of control law

				u = int(controller.Update(s.Ctler, s.ArrivalRate))

				// update PC
				s.PC = u

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
		default:
		}
	}
}

// Configure publisher RabbitMQ
func (s *Server) configureRabbitMQ() {
	err := error(nil)

	// create connection
	s.Conn, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
	//s.Conn, err = amqp.Dial("amqp://nsr:nsr@localhost:5672/") // Docker 'some-rabbit'
	//s.Conn, err = amqp.Dial("amqp://guest:guest@localhost:5672/") // Docker 'some-rabbit'
	//s.Conn, err = amqp.Dial("amqp://guest:guest@172.22.70.20:5672/") // Docker 'some-rabbit'

	shared.FailOnError(err, "Failed to connect to RabbitMQ")

	//connSub, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	//s.ConnSub, err = amqp.Dial("amqp://nsr:nsr@localhost:5672/") // Docker 'some-rabbit'
	//shared.FailOnError(err, "Failed to connect to RabbitMQ - Subscriber")
	//defer conn.Close()

	// create channel
	s.Ch, err = s.Conn.Channel()
	shared.FailOnError(err, "Failed to open a channel")

	// declare queues
	s.Queue, err = s.Ch.QueueDeclare(
		"rpc_queue", // name
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	shared.FailOnError(err, "Failed to declare a queue")

	// create a consumer
	s.Msgs, err = s.Ch.Consume(
		s.Queue.Name, // queue
		"",           // consumer
		false,        // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	shared.FailOnError(err, "Failed to register a consumer")

	// configure initial QoS of Req channel
	err = s.Ch.Qos(
		int(s.PC), // prefetch count
		0,         // prefetch size
		true,      // global TODO default is false
	)
	shared.FailOnError(err, "Failed to set QoS")
	return
}

type StatsMonitor struct {
	Alloc,
	TotalAlloc,
	Sys,
	Mallocs,
	Frees,
	LiveObjects,
	PauseTotalNs uint64

	NumGC        uint32
	NumGoroutine int
} // TODO

func Monitor(duration int, sampleSize int) {
	var m StatsMonitor
	var rtm runtime.MemStats
	var interval = time.Duration(duration) * time.Second

	time.Sleep(5 * time.Second) // take a time before starting monitoring
	count := 0

	fmt.Println("Alloc ; TotalAlloc ; Sys ; Mallocs ; Frees ; LiveObjects ; PauseTotalNs ; NumGC ; NumGoroutine")
	for {
		<-time.After(interval)

		// Read full mem stats
		runtime.ReadMemStats(&rtm)

		// Number of goroutines
		m.NumGoroutine = runtime.NumGoroutine()

		// Misc memory stats
		m.Alloc = rtm.Alloc
		m.TotalAlloc = rtm.TotalAlloc
		m.Sys = rtm.Sys
		m.Mallocs = rtm.Mallocs
		m.Frees = rtm.Frees

		// Live objects = Mallocs - Frees
		m.LiveObjects = m.Mallocs - m.Frees

		// GC Stats
		m.PauseTotalNs = rtm.PauseTotalNs
		m.NumGC = rtm.NumGC

		fmt.Println(m.Alloc, ";", m.TotalAlloc, ";", m.Sys, ";", m.Mallocs, ";", m.Frees, ";", m.LiveObjects, ";", m.PauseTotalNs, ";", m.NumGC, ";", m.NumGoroutine)

		count++
		if count > sampleSize {
			os.Exit(0)
		}
	}
} // TODO
