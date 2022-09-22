package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"rabbitmq/shared"
	"sync"
	"time"

	"github.com/streadway/amqp"
)

type Client struct {
	Id         string
	MsgSize    int
	SampleSize int
	Mean       float64
	StdDev     float64
	Conn       *amqp.Connection
	Ch         *amqp.Channel
	Queue      amqp.Queue
	Msgs       <-chan amqp.Delivery
}

func main() { // Windows
	var ws sync.WaitGroup

	numOfClients := 50 // for experimental purposes

	// configure/read flags
	var clientIdPtr = flag.String("publisher-id", "1", "publisher-id is an int")
	var msgSizePtr = flag.Int("message-size", 256, "message-size is an int")
	var sampleSizePtr = flag.Int("sample-size", 1, "sample-size is an int")
	var meanRequestTimePtr = flag.Int("mean-request-time", 1, "mean-request-time is an int (ms)")
	var stdDevMeanRequestTimePtr = flag.Int("std-dev-mean-request-time", 0, "std-dev-mean-request-time is an int")
	flag.Parse()

	// make requests to consumer
	//totalTime := c.Run()
	for { // experimental purpose
		for i := 0; i < numOfClients; i++ {
			// create publisher
			c := NewClient(*clientIdPtr, *msgSizePtr, *sampleSizePtr, *meanRequestTimePtr, *stdDevMeanRequestTimePtr)

			ws.Add(1)
			go c.RunWindows(&ws)
		}
		fmt.Println("Publishers started [", numOfClients, "publishers", *msgSizePtr, "bytes", "MeanRequestTime", *meanRequestTimePtr, "ms", "STDEV=", *stdDevMeanRequestTimePtr, "]")
		ws.Wait()
		fmt.Println("All", numOfClients, "clients finished...")
	}
}

func NewClient(clientIdPtr string, msgSizePtr int, sampleSizePtr int, meanRequestTimePtr int, stdDevMeanRequestTimePtr int) Client {
	c := Client{}

	// random setup
	rand.Seed(time.Now().UTC().UnixNano())

	// configure publisher
	c.Id = clientIdPtr
	c.MsgSize = msgSizePtr
	c.SampleSize = sampleSizePtr
	c.Mean = float64(meanRequestTimePtr)
	c.StdDev = float64(stdDevMeanRequestTimePtr)

	// Configure rabbitmq elements
	c.configureRabbitMQ()

	return c
}

func (c Client) RunWindows(ws *sync.WaitGroup) time.Duration {
	defer ws.Done()

	// Close channels and connections (when finish)
	defer c.Conn.Close()
	//defer c.RepConn.Close()
	defer c.Ch.Close()

	// initialize variables
	err := error(nil)
	//totalTime := time.Duration(0)

	// set message
	msg := make([]uint8, c.MsgSize)

	if err != nil {
		log.Fatal(err)
	}

	// make requests
	startTime := time.Now()

	//for i := 0; i < c.SampleSize; i++ {
	for { // TODO experimental purpose
		corrId := shared.RandomString(32)

		// make resquests randomly distributed -- experimental purpose -- comment
		interTime := c.Mean + rand.NormFloat64()*c.StdDev
		time.Sleep(time.Duration(interTime) * time.Millisecond)

		err = c.Ch.Publish(
			"",          // exchange
			"rpc_queue", // routing key
			false,       // mandatory
			false,       // immediate

			amqp.Publishing{
				ContentType:   "text/plain",
				CorrelationId: corrId,
				ReplyTo:       c.Queue.Name,
				Body:          msg,
				//AppId:         c.Id, // TODO - include
				//Timestamp: time.Now(), // TODO remove
			})
		shared.FailOnError(err, "Failed to publish a message")

		//fmt.Println("Client ", c.Id, " published message >> ", corrId)

		// Receive response
		/*
			for d := range c.Msgs {
				if corrId == d.CorrelationId {
					//response := d.Body // discard result
					endTime := time.Now()
					//fmt.Println("Response: ", string(response))
					totalTime += endTime.Sub(startTime)
					//fmt.Println("Client ", c.Id, " received message << ", d.CorrelationId)
					break
				}
			}
		*/
	}

	// inspect queue -- experimental purpose
	//for {
	//	q1, err1 := c.Ch.QueueInspect("rpc_queue")
	//	shared.FailOnError(err1, "Client:: Failed to inspect the queue")
	//
	//		if q1.Messages == 0 {
	//			time.Sleep(10 * time.Second)
	//			return time.Now().Sub(startTime)
	//		}
	//	}

	return time.Now().Sub(startTime)
}

func (c Client) RunMac() time.Duration {

	// Close channels and connections (when finish)
	defer c.Conn.Close()
	//defer c.RepConn.Close()
	defer c.Ch.Close()

	// initialize variables
	err := error(nil)
	//totalTime := time.Duration(0)

	// set message
	msg := make([]uint8, c.MsgSize)

	if err != nil {
		log.Fatal(err)
	}

	// make requests
	startTime := time.Now()
	for i := 0; i < c.SampleSize; i++ {
		corrId := shared.RandomString(32)

		// make resquest randomly distributed
		interTime := c.Mean + rand.NormFloat64()*c.StdDev
		time.Sleep(time.Duration(interTime) * time.Millisecond)

		err = c.Ch.Publish(
			"",          // exchange
			"rpc_queue", // routing key
			false,       // mandatory
			false,       // immediate

			amqp.Publishing{
				ContentType:   "text/plain",
				CorrelationId: corrId,
				ReplyTo:       c.Queue.Name,
				Body:          msg,
				//AppId:         c.Id, // TODO - include
				//Timestamp: time.Now(), // TODO remove
			})
		shared.FailOnError(err, "Failed to publish a message")

		//fmt.Println("Client ", c.Id, " published message >> ", corrId)

		// Receive response
		/*
			for d := range c.Msgs {
				if corrId == d.CorrelationId {
					//response := d.Body // discard result
					endTime := time.Now()
					//fmt.Println("Response: ", string(response))
					totalTime += endTime.Sub(startTime)
					//fmt.Println("Client ", c.Id, " received message << ", d.CorrelationId)
					break
				}
			}
		*/
	}
	return time.Now().Sub(startTime)
}

func (c *Client) configureRabbitMQ() {

	err := error(nil)

	// c.Conn, err = amqp.Dial("amqp://guest:guest@192.168.0.7:5672/") // Home 192
	// c.Conn, err = amqp.Dial("amqp://guest:guest@127.0.0.1:5672/") // Home 127
	c.Conn, err = amqp.Dial("amqp://guest:guest@192.168.1.158:5672/") // Home 127
	//c.Conn, err = amqp.Dial("amqp://guest:guest@172.22.38.75:5672/") // Ufpe
	shared.FailOnError(err, "Failed to connect to RabbitMQ")

	c.Ch, err = c.Conn.Channel()
	shared.FailOnError(err, "Failed to open a channel")

	// Queue - it creates a queue if it does not exist
	c.Queue, err = c.Ch.QueueDeclare(
		"rpc_queue", // name
		false,       // durable default is false
		false,       // delete when unused
		false,       // exclusive default is true
		false,       // noWait
		nil,         // arguments
	)
	shared.FailOnError(err, "Failed to declare a queue")
}

func mainMac() {

	// configure/read flags
	var clientIdPtr = flag.String("publisher-id", "1", "publisher-id is an int")
	var msgSizePtr = flag.Int("message-size", 256, "message-size is an int")
	var sampleSizePtr = flag.Int("sample-size", 1, "sample-size is an int")
	var meanRequestTimePtr = flag.Int("mean-request-time", 1, "mean-request-time is an int (ms)")
	var stdDevMeanRequestTimePtr = flag.Int("std-dev-mean-request-time", 0, "std-dev-mean-request-time is an int")
	flag.Parse()

	// create publisher
	c := NewClient(*clientIdPtr, *msgSizePtr, *sampleSizePtr, *meanRequestTimePtr, *stdDevMeanRequestTimePtr)

	// make requests to publisher
	//totalTime := c.Run()
	c.RunMac()

	// print time
	//meanTime := float64(totalTime) / 1000000.0 / float64(c.SampleSize)
	//_ = float64(totalTime) / 1000000.0 / float64(c.SampleSize)
	//fmt.Printf("Mean 'response time': %.3f (ms) \n", meanTime)
	//fmt.Printf("%.3f\n", meanTime)
	//fmt.Printf("%.3f \n", meanTime)

	//fmt.Scanln()
}
