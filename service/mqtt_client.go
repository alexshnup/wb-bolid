package service

import (
	"log"

	mqtt "github.com/alexshnup/mqtt"
	"github.com/alexshnup/uuid"
)

// NewMqttClient open mqtt connection
func NewMqttClient(protocol, address, port string, qos byte) (mqtt.Client, error) {
	// generate new uuid
	id := uuid.New()
	// concat address
	server := protocol + "://" + address + ":" + port

	clientOptions := mqtt.NewClientOptions()
	clientOptions.AddBroker(server)
	clientOptions.SetClientID(id)
	clientOptions.SetDefaultPublishHandler(defaultMessageHandler)

	// new client
	client := mqtt.NewClient(clientOptions)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return client, nil
}

func defaultMessageHandler(client mqtt.Client, msg mqtt.Message) {
	log.Printf("TOPIC: %s, MSG: %s\n", msg.Topic(), msg.Payload())
}
