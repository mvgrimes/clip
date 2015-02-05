package commands

import (
	"log"
	"os/exec"
	"time"

	"github.com/spf13/cobra" // cli
	"github.com/spf13/viper" // Config file parsing
	"github.com/streadway/amqp"
)

var pasteCmd = &cobra.Command{
	Use:   "paste",
	Short: "Watch for copied text and add it to clipboard",
	Run: func(cmd *cobra.Command, args []string) {
		InitializeConfig()
		pull()
	},
}

func connect() (*amqp.Connection, *amqp.Channel, <-chan amqp.Delivery, error) {
	logf("[*] Attempting to connect to server")

	conn, err := amqp.Dial(viper.GetString("server"))
	// failOnError(err, "Failed to connect to RabbitMQ")
	if err != nil {
		return nil, nil, nil, err
	}
	// defer conn.Close()

	ch, err := conn.Channel()
	// failOnError(err, "Failed to open a channel")
	if err != nil {
		return nil, nil, nil, err
	}
	// defer ch.Close()

	err = ch.ExchangeDeclare(
		viper.GetString("exchange"), // exchange
		"fanout",                    // type
		true,                        // durable
		false,                       // auto-deleted
		false,                       // internal
		false,                       // no-wait
		nil,                         // arguments
	)
	// failOnError(err, "Failed to declare an exchange")
	if err != nil {
		return nil, nil, nil, err
	}

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when usused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	// failOnError(err, "Failed to declare a queue")
	if err != nil {
		return nil, nil, nil, err
	}

	err = ch.QueueBind(
		q.Name, // queue name
		"",     // routing key
		viper.GetString("exchange"), // exchange
		false,
		nil)
	// failOnError(err, "Failed to bind a queue")
	if err != nil {
		return nil, nil, nil, err
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	// failOnError(err, "Failed to register a consumer")
	if err != nil {
		return nil, nil, nil, err
	}

	logf("[=] Connected to queue %s", q.Name)

	return conn, ch, msgs, err
}

func pull() {

	go func() {
		key := getKey()
		sleeper := makeSleeper()

		for {
			conn, ch, msgs, err := connect()
			// failOnError(err, "[x] error connecting")
			if err != nil {
				log.Printf("[e] error connecting: %s", err)
				sleeper()
				continue
			}
			sleeper = makeSleeper()

			defer conn.Close()
			defer ch.Close()

			for d := range msgs {
				plainText, err := decrypt(key, d.Body)

				if err != nil {
					log.Printf("[x] error in decrypt: %s", err)
					continue
				}

				if viper.GetBool("Verbose") {
					log.Printf("[>] %s", plainText)
				}

				sendToClipboard(plainText)
			}

			logf("[=] Connection closed")
		}
	}()

	logf("[*] Waiting for logs. To exit press CTRL+C")
	defer logf("[*] Exiting...")

	forever := make(chan bool)
	<-forever
}

func sendToClipboard(text []byte) {
	cmd := exec.Command(viper.GetString("PasteCmd"))
	cmdIn, _ := cmd.StdinPipe()
	err := cmd.Start()
	failOnError(err, "Unable to run PasteCmd")
	cmdIn.Write(text)
	cmdIn.Close()
	cmd.Wait()
}

func makeSleeper() func() {
	delay := 5

	return func() {
		log.Printf("[e] ... sleeping for %d seconds", delay)
		time.Sleep(time.Second * time.Duration(delay))
		log.Printf("[e] ... up")
		delay *= 2
		if delay > 60 {
			delay = 60
		}
	}
}
