package main

import (
	"fmt"
	"github.com/Omegapoint/golang-umea-20200110/Protocol"
	uuid "github.com/satori/go.uuid"
)

type ClientMap = map[uuid.UUID]Protocol.Client

func main() {
	conf := getConfig()
	fmt.Printf("Using config: %v\n", conf)

	subscription := subscribe(conf)
	fmt.Printf("successfully connected to name server at: %s:%v\n", conf.NameServerIp, conf.NameServerPort)

	for true {
		clients := <- subscription
		for _, client := range clients {
			fmt.Printf("{id: %s, ip: %s, port: %v, name: %s, connected: %v}\n", client.Id(), client.Ip().String(), client.Port(), client.Name(), client.Connected())
		}
	}
}

