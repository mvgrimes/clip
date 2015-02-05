package commands

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/cobra" // cli
	"github.com/spf13/viper" // Config file parsing
	"github.com/streadway/amqp"
)

var copyCmd = &cobra.Command{
	Use:   "copy",
	Short: "Read from STDIN and push to all paste instances",
	Run: func(cmd *cobra.Command, args []string) {
		InitializeConfig()
		push()
	},
}

func push() {
	// conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	conn, err := amqp.Dial(viper.GetString("server"))
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		viper.GetString("exchange"), // exchange
		"fanout",                    // type
		true,                        // durable
		false,                       // auto-deleted
		false,                       // internal
		false,                       // no-wait
		nil,                         // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	body := bodyFrom(os.Args)
	err = ch.Publish(
		viper.GetString("exchange"), // exchange
		"",    // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			// Body:        []byte(body),
			Body: body,
		})
	failOnError(err, "Failed to publish a message")

	if viper.GetBool("verbose") {
		log.Printf("[x] Sent %s", body)
	}
}

func bodyFrom(args []string) []byte {
	text, _ := ioutil.ReadAll(os.Stdin)

	// fmt.Println(text)
	key := getKey()
	cipherText, _ := encrypt(key, text)

	return cipherText
}
