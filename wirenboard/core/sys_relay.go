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

	topic := l.topic + "/" + deviceID + "/" + relayID + "/status/relay"

	// publish result
	if token := l.client.Publish(topic, qos, true, l.status); token.Wait() && token.Error() != nil {
		log.Println(token.Error())
	}

	// debug
	if l.debug {
		log.Println("[PUB]", qos, topic, l.status)
	}
}

// PublishStatus Relay status
func (l *relay) PublishStatusSensor(qos byte, deviceID, relayID string) {

	topic := l.topic + "/" + deviceID + "/" + relayID + "/status/sensor"

	// publish result
	if token := l.client.Publish(topic, qos, false, l.status); token.Wait() && token.Error() != nil {
		log.Println(token.Error())
	}

	// debug
	if l.debug {
		log.Println("[PUB]", qos, topic, l.status)
	}
}

// PublishStatus Relay status
func (l *relay) PublishADC(qos byte, deviceID, adcID string) {

	topic := l.topic + "/" + deviceID + "/" + adcID + "/status/adc"

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

	switch s_fields[2] {
	case "on":
		// receive message and DO
		switch string(msg.Payload()) {
		case "0":
			// logic when OFF

			l.status = RelayOnOff(uint8(device_id), uint8(relay_id), 0)
			log.Println("l.status", l.status)

			l.PublishStatus(0, s_fields[0], s_fields[1])
		case "1":
			// logic when ON

			l.status = RelayOnOff(uint8(device_id), uint8(relay_id), 1)
			log.Println("l.status", l.status)
			l.PublishStatus(0, s_fields[0], s_fields[1])
		case "3":
			// logic when ON

			l.status = RelayWhile(uint8(device_id), uint8(relay_id), 3)
			log.Println("l.status", l.status)
			l.PublishStatus(0, s_fields[0], s_fields[1])
		}

	case "sensor":
		// publish status
		l.status = Status(uint8(device_id), uint8(relay_id))
		log.Println("l.status", l.status)
		l.PublishStatusSensor(0, s_fields[0], s_fields[1])

	case "setrelaydefaultmode":
		// receive message and DO
		switch string(msg.Payload()) {
		case "off":
			// publish status
			l.status = SetConfig(uint8(device_id), uint8(relay_id), 2)
			log.Println("l.status", l.status)
			l.PublishStatus(0, s_fields[0], s_fields[1])
		case "on":
			// publish status
			l.status = SetConfig(uint8(device_id), uint8(relay_id), 1)
			log.Println("l.status", l.status)
			l.PublishStatus(0, s_fields[0], s_fields[1])
		case "blink":
			// publish status
			l.status = SetConfig(uint8(device_id), uint8(relay_id), 9)
			log.Println("l.status", l.status)
			l.PublishStatus(0, s_fields[0], s_fields[1])
		case "pcn":
			// publish status
			l.status = SetConfig(uint8(device_id), uint8(relay_id), 10)
			log.Println("l.status", l.status)
			l.PublishStatus(0, s_fields[0], s_fields[1])
		}

	case "setrelaytime":
		// receive message and DO
		v, _ := strconv.ParseUint(string(msg.Payload()), 10, 64)
		if v >= 1 && v <= 60 {
			l.status = SetConfig(uint8(device_id), uint8(relay_id)+4, uint8(v))
			log.Println("l.status", l.status)
			l.PublishStatus(0, s_fields[0], s_fields[1])
		}

	case "changeaddress":
		// receive message and DO
		var newaddr uint8 = uint8(relay_id)
		if newaddr >= 1 && newaddr <= 127 {
			l.status = ChangeAddress(uint8(device_id), newaddr)
			log.Println("l.status", l.status)
			l.PublishStatus(0, s_fields[0], s_fields[1])
		}

	case "adc":
		// publish status
		l.status = ADC(uint8(device_id), uint8(relay_id))
		log.Println("l.status", l.status)
		l.PublishADC(0, s_fields[0], s_fields[1])
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
