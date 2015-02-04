package commands

import (
	"fmt"
	"github.com/spf13/viper" // Config file parsing
	"log"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf(" [e] %s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func logf(format string, v ...interface{}) {
	if viper.GetBool("Verbose") {
		log.Printf(format, v...)
	}
}
