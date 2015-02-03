package main

import (
  "fmt"
  "log"
  "os"
  // "strings"
  "io"
  "io/ioutil"
  "crypto/aes"
  "crypto/cipher"
  "crypto/rand"
  "encoding/base64"

  "github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
  if err != nil {
    log.Fatalf("%s: %s", msg, err)
    panic(fmt.Sprintf("%s: %s", msg, err))
  }
}

func main() {
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
      Body:        body,
    })
  failOnError(err, "Failed to publish a message")

  log.Printf(" [x] Sent %s", body)
}

func bodyFrom(args []string) []byte {
  text, _ := ioutil.ReadAll(os.Stdin)

  // fmt.Println(text)
  key := getKey()
  cipherText, _ := encrypt( key, text )

  return cipherText
}

func getKey() []byte {
  key := []byte("a very very very very secret key") // 32 bytes
  return key
}

func encrypt(key, text []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    b := base64.StdEncoding.EncodeToString(text)
    ciphertext := make([]byte, aes.BlockSize+len(b))
    iv := ciphertext[:aes.BlockSize]
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        return nil, err
    }
    cfb := cipher.NewCFBEncrypter(block, iv)
    cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))
    return ciphertext, nil
}
