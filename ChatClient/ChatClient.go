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
	for true {
		<- subscription
	}
}

