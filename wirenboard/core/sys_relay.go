package syscore

import (
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	mqtt "github.com/alexshnup/mqtt"
)

/*
Struct relay provides system LED[0,1] control

Topics:
	Subscribe:
		name + "/SYSTEM/LED[0,1]/ACTION		{0, 1, STATUS}
	Publish:
		name + "/SYSTEM/LED[0,1]/STATUS		{0, 1}

Methods:
	Subscribe
	Unsubscribe
	PublishStatus

Functions:
	Set trigger to [none] when subscribe
		echo none | sudo tee /sys/class/relays/relay0/trigger
	Set trigger to [mmc0] when unsubscribe
		echo mmc0 | sudo tee /sys/class/relays/relay0/trigger
	Set brightness to 1 when ON
		echo 1 | sudo tee /sys/class/relays/relay0/brightness
	Set brightness to 0 when OFF
		echo 0 | sudo tee /sys/class/relays/relay0/brightness
	Get brightness status
		sudo cat /sys/class/relays/relay0/brightness

TODO:
[ ] catch errors in relayMessageHandler

*/
type relay struct {
	client   mqtt.Client
	debug    bool
	topic    string
	status   string
	relayID  string
	deviceID string
}

// a[len(a)-1:] last char

// newRelay return new relay object.
func newRelay(c mqtt.Client, topic string, debug bool) *relay {
	return &relay{
		client:   c,
		debug:    debug,
		topic:    topic,
		status:   "0",
		relayID:  topic[len(topic)-1:],
		deviceID: topic[len(topic)-3:],
	}
}

// Subscribe to topic
func (l *relay) Subscribe(qos byte) {

	topic := l.topic + "/#"

	log.Println("[RUN] Subscribing:", qos, topic)

	if token := l.client.Subscribe(topic, qos, l.relayMessageHandler); token.Wait() && token.Error() != nil {
		log.Println(token.Error())
	}
}

// UnSubscribe from topic
func (l *relay) UnSubscribe() {

	topic := l.topic + "/#"

	log.Println("[RUN] UnSubscribing:", topic)

	if token := l.client.Unsubscribe(topic); token.Wait() && token.Error() != nil {
		log.Println(token.Error())
	}
}

// PublishStatus Relay status
func (l *relay) PublishStatus(qos byte, deviceID, relayID string) {

	topic := l.topic + "/" + deviceID + "/" + relayID + "/status"

	// publish result
	if token := l.client.Publish(topic, qos, false, l.status); token.Wait() && token.Error() != nil {
		log.Println(token.Error())
	}

	// debug
	if l.debug {
		log.Println("[PUB]", qos, topic, l.status)
	}
}

// relayMessageHandler set Relay to ON or OFF and get STATUS
func (l *relay) relayMessageHandler(client mqtt.Client, msg mqtt.Message) {

	// debug
	if l.debug {
		log.Println("[SUB]", msg.Qos(), msg.Topic(), string(msg.Payload()))
	}

	s1 := strings.Replace(msg.Topic(), l.topic, "", 1)
	s1 = strings.Replace(s1, "/", " ", -1)
	s_fields := strings.Fields(s1)
	device_id, _ := strconv.ParseUint(s_fields[0], 10, 64)
	relay_id, _ := strconv.ParseUint(s_fields[1], 10, 64)

	// receive message and DO
	switch string(msg.Payload()) {
	case "0":
		// logic when OFF

		Relay_OFF(uint8(device_id), uint8(relay_id))
		l.status = StatusRelay(uint8(device_id), uint8(relay_id))
		log.Println("l.status", l.status)

		l.PublishStatus(0, s_fields[0], s_fields[1])
	case "1":
		// logic when ON
		Relay_ON(uint8(device_id), uint8(relay_id))
		l.status = StatusRelay(uint8(device_id), uint8(relay_id))
		log.Println("l.status", l.status)
		l.PublishStatus(0, s_fields[0], s_fields[1])
	case "status":
		// publish status
		l.status = StatusRelay(uint8(device_id), uint8(relay_id))
		log.Println("l.status", l.status)
		l.PublishStatus(0, s_fields[0], s_fields[1])
	}
}

// getBrightness
func getBrightness(relayID string) (string, error) {
	dat, err := ioutil.ReadFile("/sys/class/relays/relay" + relayID + "/brightness")
	if err != nil {
		return "", err
	}

	return strings.Trim(string(dat), "\r\n"), nil
}

// setBrightness
func setBrightness(relayID, data string) error {
	err := ioutil.WriteFile("/sys/class/relays/relay"+relayID+"/brightness", []byte(data), 0644)
	if err != nil {
		return err
	}
	return nil
}
