package syscore

import (
	mqtt "github.com/alexshnup/paho.mqtt.golang"
)

// System struct
type System struct {
	Memory *memory
	Relay  *relay
	// Relay1 *relay
	// Relay2 *relay
	// Relay3 *relay
	// Relay4 *relay
}

// NewSystem return new System object.
func NewC2000(c mqtt.Client, name string, debug bool) *System {
	return &System{
		// Memory: newMemory(c, name, debug),
		Relay: newRelay(c, name, debug),
		// Relay1: newRelay(c, name+"/3/1", debug),
		// Relay2: newRelay(c, name+"/3/2", debug),
		// Relay3: newRelay(c, name+"/3/3", debug),
		// Relay4: newRelay(c, name+"/3/4", debug),
	}
}
