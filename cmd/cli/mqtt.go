package main

import (
	"fmt"
	"time"

	gomqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	broker = "broker.mqttdashboard.com"
	port   = 1883
	topic  = "/home/smart_water_heater"
)

var messagePubHandler gomqtt.MessageHandler = func(client gomqtt.Client, msg gomqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler gomqtt.OnConnectHandler = func(client gomqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler gomqtt.ConnectionLostHandler = func(client gomqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

func connect(client gomqtt.Client) {
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

func NewClient(clientID string) gomqtt.Client {

	opts := gomqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID(clientID)
	// opts.SetUsername("emqx")
	// opts.SetPassword("public")
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := gomqtt.NewClient(opts)

	return client

}

func sub(client gomqtt.Client) {
	token := client.Subscribe(topic, 1, nil)
	token.Wait()
	fmt.Printf("Subscribed to topic %s\n", topic)
}

func publish(client gomqtt.Client, msg string) {

	token := client.Publish(topic, 0, false, msg)
	token.Wait()
	time.Sleep(time.Second)

}
