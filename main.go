package main

import
// "flag"
(
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/alexshnup/wb-bolid/conf"
	"github.com/alexshnup/wb-bolid/service"
	"github.com/alexshnup/wb-bolid/wirenboard"
)

// checkError check error
func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	fmt.Printf("%v", conf.Config.Mqtt.Address)

	// interrupt
	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)

	// open mqtt connection
	client, err := service.NewMqttClient(
		conf.Config.Mqtt.Protocol,
		conf.Config.Mqtt.Address,
		conf.Config.Mqtt.Port,
		0,
	)
	checkError(err)

	// new instance of mqtt client
	wb := wirenboard.NewWirenboard(client, conf.Config.Name, conf.Config.Debug)

	// Run publisher

	// wb.System.Memory.Publish(Config.Timeout, 0)

	// Run subscribing
	wb.System.Relay.Subscribe(2)
	// wb.System.Relay1.Subscribe(2)
	// wb.System.Relay2.Subscribe(2)
	// wb.System.Relay3.Subscribe(2)
	// wb.System.Relay4.Subscribe(2)

	// wait for terminating
	for {
		select {
		case <-interrupt:
			log.Println("Clean and terminating...")

			// Unsubscribe when terminating
			wb.System.Relay.UnSubscribe()
			// wb.System.Relay1.UnSubscribe()
			// wb.System.Relay2.UnSubscribe()
			// wb.System.Relay3.UnSubscribe()
			// wb.System.Relay4.UnSubscribe()

			// disconnecting
			client.Disconnect(250)

			os.Exit(0)
		}
	}

	// LoadJSONConfig()
	//
	// fmt.Println(ConfigJSON.Address)

	// for _, m := range ConfigJSON.Relays {
	// 	fmt.Println(m)
	// 	Relay_ON(ConfigJSON.Address, m)
	// 	StatusRelay(ConfigJSON.Address, m)
	// 	time.Sleep(100 * time.Millisecond)
	// 	Relay_OFF(ConfigJSON.Address, m)
	// 	StatusRelay(ConfigJSON.Address, m)
	// }

}
