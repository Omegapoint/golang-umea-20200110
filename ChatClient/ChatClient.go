package main

import (
	"encoding/json"
	"fmt"
	"github.com/Omegapoint/golang-umea-20200110/Protocol"
	uuid "github.com/satori/go.uuid"
	"net"
	"os"
	"strconv"
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
	tcpAddr, _ := net.ResolveTCPAddr("tcp4", ":" + strconv.FormatUint(uint64(conf.Port), 10))
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		printErrorMessage(fmt.Sprintf("failed to listen for incomming connections: %v\n", err))
		os.Exit(1)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		go handleClient(conn, clients)
	}
}

func handleClient(conn net.Conn, clients ConnectionMap) {
	defer conn.Close()

	request := make([]byte, 2048)
	var message Protocol.Message
	clientName := ""
	for {
		size, err := conn.Read(request)
		if err != nil {
			printInfoMessage(fmt.Sprintf("'%s' disconnected", clientName))
			return
		}
		err = json.Unmarshal(request[:size], &message)

		if clientName == "" {
			client := clients[message.Id]
			if client == nil {
				printInfoMessage("disconnected misbehaving client")
				return
			}
			clientName = client.Name
		}

		printChatMessage(message.Message, clientName, false)
	}
}


// handleClientBroadcast is responsible for broadcasting messages received from the user to
// all the connected clients. The messages from the user is received through `userMessages`
// and `connections` contains all the active connections to outgoing clients
func handleBroadcast(clientId uuid.UUID, userMessages chan string, connections ConnectionMap) {
	for {
		msg := <- userMessages
		message, _ := Protocol.NewMessage(clientId, msg).Serialize()
		for _, conn := range connections {
			_, _ = conn.Conn.Write(message)
		}
	}
}

// manageClientConnections is responsible for managing outgoing connections to other chat clients
// the connected clients is received from the name server through `subscription` and the active
// connections will be stored in `connections`.
func manageClientConnections(subscription chan ClientMap, connections ConnectionMap) {
	for {
		clients := <-subscription
		for id, client := range clients {
			if connections[id] != nil {
				continue
			}

			ip := client.Ip().String() + ":" + strconv.FormatUint(uint64(client.Port()), 10)
			addr, err := net.ResolveTCPAddr("tcp4", ip)
			if err != nil {
				printErrorMessage(fmt.Sprintf("unable to resolve tcp address: %v\n", err))
				continue
			}

			conn, err := net.DialTCP("tcp", nil, addr)
			if err != nil {
				printErrorMessage(fmt.Sprintf("unable to establish connection: %v\n", err))
				continue
			}

			printInfoMessage(fmt.Sprintf("'%s' connected", client.Name()))

			connection := clientConnection{
				Name: client.Name(),
				Conn: conn,
			}
			connections[id] = &connection
		}
	}
}


