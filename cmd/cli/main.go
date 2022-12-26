package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var (
	isConnected = false
	maxTemp     = 55
)

// will use to access the container at the indexes.
const (
	// iota starts at 0
	labelIndex = iota
	inputIndex
)

func main() {

	if err := run(); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

}

func run() error {

	client := NewClient("clientId-3UD11rGWHU")

	a := app.New()
	a.Settings().SetTheme(theme.DarkTheme())
	w := a.NewWindow("Smart water heater")
	// fix starting window size
	w.Resize(fyne.NewSize(400, 400))

	btn := connectBtn(client, w)

	tempCont := newContainer("Current Temperature:", startingTemp(), false)
	topicCont := newContainer("Pub topic:", topic, false)
	maxTempCont := newContainer("Desired temperature:", fmt.Sprintf("%d", maxTemp), false)

	go UpdateTemp(client, tempCont.Objects[inputIndex])

	mainContainer := container.NewVBox(tempCont, topicCont, maxTempCont, btn)
	w.SetContent(mainContainer)
	// Run the app
	w.ShowAndRun()

	return nil
}

// UpdateTemp will simulate getting info from mqtt and changing the temperature on screen.
func UpdateTemp(client mqtt.Client, input fyne.CanvasObject) {
	// Create a timer that fires every 3 seconds
	ticker := time.NewTicker(3 * time.Second)

	// Start a loop to run the function every time the timer fires.
	for range ticker.C {

		// only update if the client is connected.
		if isConnected {
			inputVal, _ := strconv.Atoi(input.(*widget.Entry).Text)

			// if we got to the max temperature than don't update.
			if maxTemp == inputVal {
				publish(client, fmt.Sprintf("Reached desired temperature: %dc", maxTemp))
				continue
			}

			inputVal += 1
			input.(*widget.Entry).Text = fmt.Sprintf("%d", inputVal)

			// send the msg to the mqtt client.
			msg := fmt.Sprintf("Current Temperature: %dc", inputVal)
			publish(client, msg)

			// refresh the input so we display the updated value.
			input.(*widget.Entry).Refresh()

		}
	}

}

// newContainer creates a new container with the label and the input inside.
func newContainer(labelText string, inputValue string, enable bool) *fyne.Container {
	// Create an input field and a label
	input := widget.NewEntry()
	if !enable {
		input.Disable()
	}

	input.SetText(inputValue)
	label := widget.NewLabel(labelText)

	// Create a container and add the input field and the label as children
	container := container.NewVBox(label, input)
	input.Resize(label.Size())
	container.Resize(fyne.NewSize(500, 25))
	return container
}

// connectBtn returns the connection button
func connectBtn(client mqtt.Client, w fyne.Window) *widget.Button {
	btn := widget.NewButton("Connect", nil)
	btn.OnTapped = func() {
		if !isConnected {
			connect(client)
			sub(client)
			// publish(client)
			isConnected = true
			btn.SetText("Disconnect")
		} else {
			client.Disconnect(25)
			btn.SetText("Connect")

			isConnected = false

		}

	}

	return btn
}

// startingTemp will return a random staring temperature.
func startingTemp() string {
	// Seed the random number generator with the current time
	rand.Seed(time.Now().UnixNano())

	// Generate a random integer between 1 and 10
	randomInt := rand.Intn(10) + 1

	return fmt.Sprintf("%d", randomInt)
}
