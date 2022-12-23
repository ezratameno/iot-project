package main

import (
	"fmt"
	"os"

	mqtt "github.com/ezratameno/iot-project/internal/mqtt"
)

func main() {

	if err := run(); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

}

func run() error {
	mqtt.Start("clientId-41HiffZGq8")
	return nil
}
