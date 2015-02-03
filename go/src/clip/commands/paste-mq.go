package commands

import (
	// "fmt"
	"log"
	"os/exec"

	"github.com/spf13/cobra" // cli
	"github.com/streadway/amqp"
)

var pasteCmd = &cobra.Command{
	Use:   "paste",
	Short: "Watch for copied text and add it to clipboard",
	Run: func(cmd *cobra.Command, args []string) {
		pull()
	},
}

func pull() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs",   // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when usused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name, // queue name
		"",     // routing key
		"logs", // exchange
		false,
		nil)
	failOnError(err, "Failed to bind a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			key := getKey()
			plainText, _ := decrypt(key, d.Body)

			log.Printf(" [x] %s", plainText)

			cmd := exec.Command("pbcopy")
			cmdIn, _ := cmd.StdinPipe()
			cmd.Start()
			cmdIn.Write(plainText)
			cmdIn.Close()
			cmd.Wait()
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}
