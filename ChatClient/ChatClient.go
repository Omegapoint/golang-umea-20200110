package main

import (
	"fmt"
	"github.com/Omegapoint/golang-umea-20200110/Protocol"
	uuid "github.com/satori/go.uuid"
	"net"
)

type ClientMap = map[uuid.UUID]Protocol.Client
type clientConnection struct {
	Name string
	Conn *net.TCPConn
}
type ConnectionMap = map[uuid.UUID]*clientConnection

func main() {
	printJumboMessage("GoChat")

	conf := getConfig()
	fmt.Println()
	printInfoMessage(fmt.Sprintf("Using config: %v", conf))

	subscription, clientId := subscribeToNameServer(conf)
	printInfoMessage(fmt.Sprintf("successfully connected to name server at: %s:%v\n", conf.NameServerIp, conf.NameServerPort))

	userMessages := make(chan string)
	go userMessageRPLoop(userMessages)

	connections := make(ConnectionMap)
	go manageClientConnections(subscription, connections)
	go handleBroadcast(clientId, userMessages, connections)
	receiveMessages(conf, connections)
}

// receiveMessages handles the receiving of messages from other clients and writing them to the user.
// `config` contains the config information that this client presented to the name server, i.e. the info
// that other clients will use when trying to connect to this client. The argument `clients` contains
// all outgoing connections from this client.
func receiveMessages(conf config, clients ConnectionMap) {
	fmt.Println("unimplemented method: receiveMessage")
}


// handleClientBroadcast is responsible for broadcasting messages received from the user to
// all the connected clients. The messages from the user is received through `userMessages`
// and `connections` contains all the active connections to outgoing clients
func handleBroadcast(clientId uuid.UUID, userMessages chan string, connections ConnectionMap) {
	fmt.Println("unimplemented method: handleBroadcast")
}

// manageClientConnections is responsible for managing outgoing connections to other chat clients
// the connected clients is received from the name server through `subscription` and the active
// connections will be stored in `connections`.
func manageClientConnections(subscription chan ClientMap, connections ConnectionMap) {
	fmt.Println("unimplemented method: manageClientConnections")
}


