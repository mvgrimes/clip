package commands

import (
	// "fmt"
	"log"
	"os"
	// "strings"
	"io/ioutil"

	"github.com/spf13/cobra" // cli
	"github.com/streadway/amqp"
)

var copyCmd = &cobra.Command{
	Use:   "copy",
	Short: "Read from STDIN and push to all paste instances",
	Run: func(cmd *cobra.Command, args []string) {
		push()
	},
}

func push() {
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

	body := bodyFrom(os.Args)
	err = ch.Publish(
		"logs", // exchange
		"",     // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			// Body:        []byte(body),
			Body: body,
		})
	failOnError(err, "Failed to publish a message")

	log.Printf(" [x] Sent %s", body)
}

func bodyFrom(args []string) []byte {
	text, _ := ioutil.ReadAll(os.Stdin)

	// fmt.Println(text)
	key := getKey()
	cipherText, _ := encrypt(key, text)

	return cipherText
}
