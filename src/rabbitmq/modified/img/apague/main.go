package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	_ "net/http/pprof"
	"os"
	"rabbitmq/shared"
)

func main() {
	go Server()
	go Client()

	fmt.Scanln()
}

func Server() {

	// create connection
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/") // Docker 'some-rabbit'
	shared.FailOnError(err, "Failed to connect to RabbitMQ")

	// create channel
	ch, err := conn.Channel()
	shared.FailOnError(err, "Failed to open a channel")

	// declare queue
	queue, err := ch.QueueDeclare(
		"rpc_queue", // name
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	shared.FailOnError(err, "Failed to declare Req queue")

	// create consumer
	msgs, err := ch.Consume(
		queue.Name, // request queue
		"",         // consumer
		false,      // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)

	for d := range msgs {
		// send ack to broker as soon the message has been received
		d.Ack(false)

		// unmarshall message
		fileContent := d.Body

		// send file back to publisher
		response := fileContent // TODO - just reply the file content

		// publish response
		err := ch.Publish(
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
	}
}

func Client() {
	err := error(nil)

	// create connection
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/") // Docker
	shared.FailOnError(err, "Failed to connect to RabbitMQ")

	// create a channel
	ch, err := conn.Channel()
	shared.FailOnError(err, "Failed to open a channel")

	// Close channels and connections (when finish)
	defer conn.Close()
	defer ch.Close()

	// create a queue if it does not exist
	queue, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // noWait
		nil,   // arguments
	)
	shared.FailOnError(err, "Client failed to declare a Request queue")

	// create a consumer
	msgs, err := ch.Consume(
		queue.Name, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	shared.FailOnError(err, "Failed to register a consumer")

	// read file
	fileContent, err := os.ReadFile("/Volumes/GoogleDrive/Meu Drive/go/adaptive/src/rabbitmq/modified/pubsub/publisher/" + "file-199350.txt")
	if err != nil {
		log.Fatal(err)
	}

	// make requests
	for i := 0; i < 100; i++ {
		corrId := shared.RandomString(32)

		err = ch.Publish(
			"",          // exchange
			"rpc_queue", // routing key
			false,       // mandatory
			false,       // immediate

			amqp.Publishing{
				ContentType:   "text/plain",
				CorrelationId: corrId,
				ReplyTo:       queue.Name,
				Body:          fileContent,
				//AppId:         c.Id, // TODO - include
				//Timestamp: time.Now(), // TODO remove
			})
		shared.FailOnError(err, "Failed to publish a message")

		// Receive response
		for d := range msgs {
			if corrId == d.CorrelationId {
				fmt.Println("File received")
				break
			}
		}
	}
	fmt.Scanln()
}
