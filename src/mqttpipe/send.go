package mqttpipe

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Message struct {
	Topic	string
	Payload	string
}

var (
	Send chan Message
)

func Sender(client mqtt.Client) {
	for {
		msg := <- Send
		token := client.Publish(msg.Topic,0, false,[]byte(msg.Payload))
		token.Wait()
	}
}