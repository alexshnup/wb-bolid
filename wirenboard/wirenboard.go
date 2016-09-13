package wirenboard

import (
	mqtt "github.com/alexshnup/mqtt"

	"github.com/alexshnup/wb-bolid/wirenboard/core"
)

// Wirenboard struct
type Wirenboard struct {
	client mqtt.Client
	name   string
	System *syscore.System
}

func NewWirenboard(c mqtt.Client, name string, debug bool) *Wirenboard {
	return &Wirenboard{
		client: c,
		name:   name,
		System: syscore.NewC2000(c, name, debug),
	}
}
