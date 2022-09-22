package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand"
	_ "net/http/pprof"
	"os"
	"rabbitmq/acmsac/controller"
	"rabbitmq/acmsac/monitor"
	"rabbitmq/shared"
	"runtime"
	"strconv"
	"time"

	"github.com/streadway/amqp"
)

var kpPtr = flag.Float64("kp", 1.0, "Kp is a float")
var kiPtr = flag.Float64("ki", 1.0, "Ki is a float")
var kdPtr = flag.Float64("kd", 1.0, "Kd is a float")
var csvLabel = flag.String("label", "training", "Label is a String")
var samples = 100 - 1 // Minus due to count starting from 0

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
	csvPrinter      bool
}

func main() {

	// configure/read flags
	var isAdaptivePtr = flag.Bool("is-adaptive", false, "is-adaptive is a boolean")
	var controllerTypePtr = flag.String("controller-type", "OnOff", "controller-type is a string")
	var monitorIntervalPtr = flag.Int("monitor-interval", 1, "monitor-interval is an int (s)")
	var setPointPtr = flag.Float64("set-point", 3000.0, "set-point is a float (goal rate)")
	var prefetchCountPtr = flag.Int("prefetch-count", 1, "prefetch-count is an int")
	var csvPrinter = flag.Bool("csv-printer", false, "csv-printer is a bool")
	flag.Parse()

	// create controller
	var c controller.IController
	c = controller.NewController(*controllerTypePtr, *setPointPtr, *kpPtr, *kiPtr, *kdPtr)

	// create new consumer
	var server = NewServer(*isAdaptivePtr, *monitorIntervalPtr, c, *prefetchCountPtr, *csvPrinter)

	fmt.Println("Consumer started [", *isAdaptivePtr, *controllerTypePtr, "Kp=", *kpPtr, "Ki=", *kiPtr, "Kd=", *kdPtr, "Goal=", *setPointPtr, "Monitor Interval=", *monitorIntervalPtr, "PC=", *prefetchCountPtr, "]")

	// run consumer
	server.Run()
}

func NewServer(isAdaptive bool, monitorInterval int, c controller.IController, prefetchCount int, csvPrinter bool) Server {
	s := Server{}

	// Configure consumer
	s.IsAdaptive = isAdaptive
	s.MonitorInterval = time.Duration(monitorInterval) * time.Second

	// Initialise channel to communicate with Monitor
	s.ChStart = make(chan bool)
	s.ChStop = make(chan bool)

	// create Monitor
	s.Mnt = monitor.NewMonitor(s.MonitorInterval)

	// set controller
	s.Ctler = c

	// set initial PC -- always 1
	s.PC = prefetchCount
	s.csvPrinter = csvPrinter
	return s
}

// Run consumer
func (s Server) Run() {

	// close all rabbitmq elements before exiting
	defer s.Conn.Close()
	defer s.Ch.Close()

	// start timer
	go s.Mnt.Monitoring(s.ChStart, s.ChStop)

	// Configure RabbitMQ
	s.configureRabbitMQ()

	// handle requests
	s.handleRequests()
	//s.handleRequestsZieglerNicholsTraining()
}

// Handle requests
func (s Server) handleRequests() {

	// initial prefetch
	u := s.PC
	count := 0       // count received messages
	countSample := 0 // for experimental purpose
	t1 := time.Time{}
	csvFile := &os.File{}
	err := error(nil)
	csv_writer := csv.NewWriter(nil)
	if s.csvPrinter { //Create file if not addaptative
		currentTime := time.Now()
		kpStr := fmt.Sprintf("%.4f", *kpPtr)
		kdStr := fmt.Sprintf("%.4f", *kdPtr)
		kiStr := fmt.Sprintf("%.4f", *kiPtr)
		csvFile, err = os.Create("../sheets/" + *csvLabel + "-" + currentTime.Format("02-01-2006") + "-kp-" + kpStr + "-kd-" + kdStr + "-ki-" + kiStr + ".csv")
		if err != nil {
			log.Fatalf("failed creating file: %s", err)
		}
		csv_writer = csv.NewWriter(csvFile)
	}
	for {
		select {
		case d := <-s.Msgs: // receive a message
			// ack message as soon as it arrives
			//time.Sleep(1 * time.Millisecond) // TODO
			// time.Sleep(250 * time.Microsecond) // TODO
			d.Ack(false)
			count++ // increment number of received messages
		case <-s.ChStart: // start timer
			t1 = time.Now()
			count = 0
		case <-s.ChStop: // stop timer

			// calculate arrival rate
			monitorInterval := time.Now().Sub(t1).Seconds()
			s.ArrivalRate = float64(count) / monitorInterval

			// inspect queue
			q1, err1 := s.Ch.QueueInspect("rpc_queue")
			shared.FailOnError(err1, "Failed to inspect the queue")

			// log information
			// FON
			formatedArrival := fmt.Sprintf("%.3f", s.ArrivalRate)
			fmt.Printf("%d;%d;%s \n", s.PC, q1.Messages, formatedArrival)
			if q1.Messages == 0 && s.ArrivalRate == 0 && s.PC == 1 {
				fmt.Println("Not yet started")
			} else if s.csvPrinter {
				csv_writer.Write([]string{strconv.Itoa(s.PC), strconv.Itoa(q1.Messages), formatedArrival})

				if countSample%samples == 0 && countSample != 0 {
					csv_writer.Flush()
					if err := csv_writer.Error(); err != nil {
						log.Fatal(err) // write file.csv: bad file descriptor
					} else {
						fmt.Println("Flushed into csv")
					}
				}
			}
			if q1.Messages < int(s.ArrivalRate) {
				fmt.Println("Stopped for 10 min")
				time.Sleep(10 * time.Minute)
			}
			if countSample < samples {
				if s.ArrivalRate != 0 {
					countSample++
				}
			} else {
				countSample = 0
				if !s.IsAdaptive {
					s.PC = s.PC + 1

					// set qos
					err := s.Ch.Qos(
						s.PC, // update prefetch count
						0,    // prefetch size
						true, // default is false
					)
					shared.FailOnError(err, "Failed to set QoS")
				}
			}
			if s.IsAdaptive { // adaptive

				// compute new value of control input using a given controller
				u = int(math.Round(controller.Update(s.Ctler, s.ArrivalRate)))

				// update PC
				s.PC = u

				// Reconfigure QoS (Ineffective if autoAck is true)
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

// Handle requests
func (s *Server) handleRequestsZieglerNicholsTraining() {

	// initial prefetch
	//u := s.PC
	count := 0       // count received messages
	countSample := 0 // for experimental purpose
	t1 := time.Time{}

	for {
		select {
		case d := <-s.Msgs: // receive a message
			// ack message as soon as it arrives
			//time.Sleep(1 * time.Millisecond) // TODO
			// time.Sleep(250 * time.Microsecond) // TODO
			d.Ack(false)
			count++ // increment number of received messages
		case <-s.ChStart: // start timer
			t1 = time.Now()
			count = 0
		case <-s.ChStop: // stop timer

			// calculate arrival rate
			monitorInterval := time.Now().Sub(t1).Seconds()
			s.ArrivalRate = float64(count) / monitorInterval

			countSample++

			// inspect queue
			q1, err1 := s.Ch.QueueInspect("rpc_queue")
			shared.FailOnError(err1, "Failed to inspect the queue")

			// log information
			fmt.Printf("%d;%d;%.3f \n", s.PC, q1.Messages, s.ArrivalRate)

			if s.PC == 1 && countSample > 30 {
				s.MonitorInterval = 1
				s.PC = 2 // From PC = 1 to 2 (step)

				// set qos
				err := s.Ch.Qos(
					s.PC, // update prefetch count
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
	// s.Conn, err = amqp.Dial("amqp://guest:guest@192.168.0.7:5672/") // Home 192
	s.Conn, err = amqp.Dial("amqp://guest:guest@127.0.0.1:5672/") // Home 192

	//s.Conn, err = amqp.Dial("amqp://guest:guest@172.22.38.75:5672/") // Home

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

func (s *Server) processRequest() {

	// Do something
	interTime := 100 + rand.NormFloat64()*10
	time.Sleep(time.Duration(interTime) * time.Millisecond)
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
