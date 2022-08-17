package main

import (
	"flag"
	"fmt"
	"github.com/streadway/amqp"
	_ "net/http/pprof"
	"os"
	"rabbitmq/modified/controller"
	"rabbitmq/modified/controller/pid"
	"rabbitmq/modified/monitor"
	"rabbitmq/shared"
	"runtime"
	"time"
)

type Server struct {
	IsAdaptive bool
	//Cnt             controller.Controller
	MonitorInterval time.Duration
	Conn            *amqp.Connection
	Ch              *amqp.Channel
	Queue           amqp.Queue
	Msgs            <-chan amqp.Delivery
	ChStart         chan bool
	ChStop          chan bool
	Mnt             monitor.Monitor
	Ctler           controller.Controller
}

func main() {

	// configure/read flags
	var isAdaptivePtr = flag.Bool("is-adaptive", false, "is-adaptive is a boolean")
	var controllerTypePtr = flag.String("controller-type", "PID", "controller-type is a string")
	var prefetchCountInitialPtr = flag.Int("prefetch-count-initial", 1, "prefetch-count-initial is an int")
	var monitorIntervalPtr = flag.Int("monitor-interval", 1, "monitor-interval is an int (ms)")
	var setPoint = flag.Float64("set-point", 1601.0, "set-point is a float (goal rate)")
	var kp = flag.Float64("kp", 1601.0, "kp is a float (constant K of PID)")
	flag.Parse()

	// create new consumer-old
	var server = NewServer(*isAdaptivePtr, *controllerTypePtr, *prefetchCountInitialPtr, *monitorIntervalPtr, *setPoint, *kp)

	// execute_consumer consumer-old
	//fmt.Println("Server is running ...")

	// start monitoring
	//go Monitor(5, 90) // seconds & 30 samples

	// run consumer-old
	server.Run()
}

func NewServer(isAdaptive bool, controllerType string, prefetchCountInitial int, monitorInterval int, setPoint float64, kp float64) Server {
	s := Server{}

	// Configure consumer-old
	s.IsAdaptive = isAdaptive
	s.MonitorInterval = time.Duration(monitorInterval) * time.Second

	// Initialise channel to communicate with Monitor
	s.ChStart = make(chan bool)
	s.ChStop = make(chan bool)

	// create Monitor
	s.Mnt = monitor.NewMonitor(s.MonitorInterval)

	// create controller
	s.Ctler = controller.NewController(controllerType, s.Mnt, prefetchCountInitial, setPoint, kp)
	s.MonitorInterval = time.Duration(monitorInterval) * time.Millisecond

	return s
}

// Run consumer-old
func (s Server) Run() {

	// close all rabbitmq elements before exiting
	defer s.Conn.Close()
	defer s.Ch.Close()

	// start monitor
	go s.Mnt.Monitoring(s.ChStart, s.ChStop)

	// Configure RabbitMQ
	s.configureRabbitMQ()

	// handle requests
	s.handlRequests()
}

// Handle requests
func (s Server) handlRequests() {
	forever := make(chan bool)

	var m runtime.MemStats // TODO

	go func(chStart, chStop chan bool) {
		count := 0

		pid := pid.PIDController{}     // TODO
		pid.Init(s.Ctler.KP, 0.0, 0.0) // TODO

		for {
		myLoop:
			for d := range s.Msgs {
				q1, err1 := s.Ch.QueueInspect("rpc_queue")
				shared.FailOnError(err1, "Failed to publish a message")

				fmt.Println(q1.Messages, time.Now().Sub(d.Timestamp).Milliseconds())

				// send ack to broker as soon the message has been received
				d.Ack(false)

				// unmarshall message
				fileContent := d.Body

				// invoke fibonacci
				response := fileContent // TODO - just reply the file content

				// publish response
				err := s.Ch.Publish(
					"",        // exchange
					d.ReplyTo, // routing key
					false,     // mandatory
					false,     // immediate
					amqp.Publishing{
						ContentType:   "text/plain",
						CorrelationId: d.CorrelationId,
						Body:          response,
					})
				shared.FailOnError(err, "Failed to publish a message")

				// interact with Monitor/Controller
				select {
				case <-chStart: // start monitor
					count = 0
				case <-chStop: // stop monitor
					s.Ctler.ProcRate = float64(count) / float64(s.Mnt.MonitorInterval)
					runtime.ReadMemStats(&m) // TODO

					// Reconfigure QoS (Ineffective if autoAck is true)
					if s.IsAdaptive {
						s.Ctler.PC = int(pid.Update(float64(s.Ctler.SP), s.Ctler.ProcRate)) // TODO
						err := s.Ch.Qos(
							s.Ctler.PC, // update prefetch count
							0,          // prefetch size
							false,      // default is false
						)
						shared.FailOnError(err, "Failed to set QoS")
					}
					break myLoop
				default: // normal processing
					count++
				}
			}
		}
	}(s.ChStart, s.ChStop)
	<-forever
}

// Configure publisher RabbitMQ (consumer-old-side)
func (s *Server) configureRabbitMQ() {
	err := error(nil)

	///connPub, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	//s.Conn, err = amqp.Dial("amqp://nsr:nsr@localhost:5672/") // Docker 'some-rabbit'
	s.Conn, err = amqp.Dial("amqp://guest:guest@localhost:5672/") // Docker 'some-rabbit'
	shared.FailOnError(err, "Failed to connect to RabbitMQ")

	//connSub, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	//s.ConnSub, err = amqp.Dial("amqp://nsr:nsr@localhost:5672/") // Docker 'some-rabbit'
	//shared.FailOnError(err, "Failed to connect to RabbitMQ - Subscriber")
	//defer conn.Close()

	s.Ch, err = s.Conn.Channel()
	shared.FailOnError(err, "Failed to open a channel")

	s.Queue, err = s.Ch.QueueDeclare(
		"rpc_queue", // name
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	shared.FailOnError(err, "Failed to declare Req queue")

	/*
		if s.Ctler.PC != 0 { // start having a infinite prefetch buffer
			err = s.ChSub.Qos(
				s.Ctler.PC,     // prefetch count
				0,   // prefetch size
				true,    // global - default false
			)
			shared.FailOnError(err, "Failed to set QoS")
		}
	*/

	s.Msgs, err = s.Ch.Consume(
		s.Queue.Name, // request queue
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
		s.Ctler.PC, // prefetch count
		0,          // prefetch size
		true,       // global TODO default is false
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
