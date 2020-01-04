package main

import (
	"bufio"
	"fmt"
	"github.com/Omegapoint/golang-umea-20200110/Protocol"
	uuid "github.com/satori/go.uuid"
	"net"
	"os"
	"strconv"
	"time"
)

type ClientMap = map[uuid.UUID]Protocol.Client
type ConnectionMap = map[uuid.UUID]*net.TCPConn

func main() {
	conf := getConfig()
	fmt.Printf("Using config: %v\n", conf)

	subscription, clientId := subscribe(conf)
	fmt.Printf("successfully connected to name server at: %s:%v\n", conf.NameServerIp, conf.NameServerPort)

	clientMessages := make(chan string)
	go userMessageRPLoop(clientMessages)
	handleClientBroadcast(clientId, subscription, clientMessages)
}

// userMessageRPLoop reads messages from the user and prints them to stdout as well as writing
// them to the provided channel.
func userMessageRPLoop(messages chan string) {
	fmt.Printf("\n\n=============================================\n")
	reader := bufio.NewReader(os.Stdin)
	var message string
	for true {
		fmt.Print("Say something: ")
		message, _ = reader.ReadString('\n')
		fmt.Printf("\033[1A")
		fmt.Printf("\033[K")
		fmt.Printf("\033[0;37m[%v]\033[0m \033[0;31m%s\033[0m: %s\n", time.Now().Format("2006-01-02 15:04:05"), "me", message)
		messages <- message
	}
}

// handleClientBroadcast is responsible for managing outgoing connections to other clients
// it receives connected clients from the name server through the channel `subscription`
// and all messages entered by the user will be received through the channel `userMessages`.
// When messages are received from the channel, this method is responsible for broadcasting
// these to all the outgoing clients.
func handleClientBroadcast(clientId uuid.UUID, subscription chan ClientMap, userMessages chan string) {
	connections := make(ConnectionMap)

	go manageClientConnections(subscription, connections)
	go handleBroadcast(clientId, userMessages, connections)
}

func handleBroadcast(clientId uuid.UUID, userMessages chan string, connections ConnectionMap) {
	for true {
		msg := <- userMessages
		message, _ := Protocol.NewMessage(clientId, msg).Serialize()
		for _, conn := range connections {
			_, _ = conn.Write(message)
		}
	}
}

func manageClientConnections(subscription chan ClientMap, connections ConnectionMap) {
	for true {
		clients := <-subscription
		for id, client := range clients {
			if connections[id] == nil {
				continue
			}

			ip := client.Ip().String() + ":" + strconv.FormatUint(uint64(client.Port()), 10)
			addr, err := net.ResolveTCPAddr("tcp4", ip)
			if err != nil {
				fmt.Fprintf(os.Stderr, "unable to resolve tcp address: %v\n", err)
				continue
			}

			conn, err := net.DialTCP("tcp", nil, addr)
			if err != nil {
				fmt.Fprintf(os.Stderr, "unable to establish connection: %v\n", err)
				continue
			}

			fmt.Printf("[%v] client '%s' connected", time.Now(), client.Name())
			connections[id] = conn
		}
	}
}


