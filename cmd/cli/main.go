package main

import (
	"fmt"
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

// b
var (
	isConnected = false
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
	w.Resize(fyne.NewSize(1200, 800))

	btn := connectBtn(client, w)

	tempCont := newContainer("Temperature:", "0")
	go OnUpdate(tempCont.Objects[1])

	mainContainer := container.NewVBox(tempCont, btn)
	w.SetContent(mainContainer)
	// Run the app
	w.ShowAndRun()

	return nil
}

// OnUpdate will simulate getting info from the mqtt and changing the screen.
func OnUpdate(input fyne.CanvasObject) {
	// Create a timer that fires every 3 seconds
	ticker := time.NewTicker(3 * time.Second)

	// Start a loop to run the function every time the timer fires
	for range ticker.C {
		intVal, _ := strconv.Atoi(input.(*widget.Entry).Text)
		intVal += 1
		input.(*widget.Entry).Text = fmt.Sprintf("%d", intVal)
		// refresh the input so we display the updated value.
		input.(*widget.Entry).Refresh()
	}
}

// newContainer creates a new container with the label and the input inside.
func newContainer(labelText, inputValue string) *fyne.Container {
	// Create an input field and a label
	input := widget.NewEntry()
	input.Disable()
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
			publish(client)
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
