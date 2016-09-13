// Package conf for go-config-manager
package conf

import (
	"log"
	"time"

	yamlCfg "github.com/alexshnup/go-config-manager/yaml"
)

// Config struct
type ConfigRoot struct {
	Debug   bool   `yaml:"debug"`
	Timeout int    `yaml:"timeout"`
	Name    string `yaml:"name"`
	Mqtt    Mqtt   `yaml:"mqtt"`
	Serial  Serial `yaml:"serial"`
	Bolid   Bolid  `yaml:"bolid"`
}

// Mqtt struct
type Mqtt struct {
	Protocol string `yaml:"protocol"`
	Address  string `yaml:"address"`
	Port     string `yaml:"port"`
}

// Serial struct
type Serial struct {
	Port        string        `yaml:"port"`
	Baud        int           `yaml:"baud"`
	ReadTimeout time.Duration `yaml:"readtimeout"`
}

// Serial struct
type Bolid struct {
	RelayON  []byte `yaml:"relayon"`
	RelayOFF []byte `yaml:"relayoff"`
}

// checkError check error
func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Config object
var Config ConfigRoot

func init() {
	// Config manager
	err := yamlCfg.NewConfig("wb-bolid-conf.yaml").Load(&Config)
	err = yamlCfg.NewConfig("/etc/wb-bolid-conf.yaml").Load(&Config)
	checkError(err)
}
