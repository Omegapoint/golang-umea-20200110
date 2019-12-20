package main

import (
	"encoding/json"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"sync"
	"time"
)

const PORT = ":8181"
const FIVE_MINUTES = time.Minute * 5

type newClientRequest struct {
	Ip string
	Port uint16
	Name string
}

var mu sync.Mutex
var connectedClients map[uuid.UUID]*Client

func main() {
	connectedClients = make(map[uuid.UUID]*Client)

	http.HandleFunc("/clients", returnConnectedClients())
	http.HandleFunc("/client", addNewClient())

	go cleanup()

	fmt.Printf("\n Server listening on: http://127.0.0.1%v \n", PORT)

	_ = http.ListenAndServe(PORT, nil)
}

func cleanup() {
	var keysToDelete []uuid.UUID

	for {
		time.Sleep(FIVE_MINUTES)
		currentTime := time.Now()
		for  key, client := range connectedClients {
			if currentTime.After(client.Connected().Add(FIVE_MINUTES)) {
				keysToDelete = append(keysToDelete, key)
			}
		}
		mu.Lock()
		for _, key  := range keysToDelete {
			delete(connectedClients, key)
		}
		mu.Unlock()
	}
}

func addNewClient() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "unsupported method", 404)
		}

		clientRequest := newClientRequest{}
		err := json.NewDecoder(r.Body).Decode(&clientRequest)
		if err != nil {
			http.Error(w, "failed to parse request body", 400)
		}

		client, err := newClient(clientRequest.Ip, clientRequest.Port, clientRequest.Name, time.Now())
		if err != nil {
			http.Error(w, err.Error(), 400)
		}
		if client == nil {
			http.Error(w, "failed to parse request body", 400)
		}

		connectedClients[client.Id()] = client
		writeConnectedClients(w)
	}
}

func returnConnectedClients() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "unsupported method", 404)
		}

		writeConnectedClients(w)
	}
}

func writeConnectedClients(w http.ResponseWriter) {
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusOK)

	responseData, err := json.Marshal(connectedClients)
	if err != nil {
		http.Error(w, "failed to retrieve response data", 500)
	}

	//Write json response back to response
	w.Write(responseData)
}
