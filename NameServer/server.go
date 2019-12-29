package main

import (
	"encoding/json"
	"fmt"
	"github.com/Omegapoint/golang-umea-20200110/Protocol"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"os"
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

type updateClientRequest struct {
	Id uuid.UUID `json:"id"`
}

var mu sync.Mutex
var connectedClients map[uuid.UUID]*Protocol.Client

func main() {
	connectedClients = make(map[uuid.UUID]*Protocol.Client)

	http.HandleFunc("/clients", returnConnectedClients())
	http.HandleFunc("/client", client)

	go cleanup()

	fmt.Printf("\n Server listening on: http://127.0.0.1%v \n", PORT)

	_ = http.ListenAndServe(PORT, nil)
}

func cleanup() {

	for {
		time.Sleep(FIVE_MINUTES)

		var keysToDelete []uuid.UUID
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

func client(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addNewClient(w, r)
	case http.MethodPatch:
		updateClient(w, r)
	default:
		http.Error(w, "usupported method", 404)
	}
}

func addNewClient(w http.ResponseWriter, r *http.Request) {
	clientRequest := newClientRequest{}
	err := json.NewDecoder(r.Body).Decode(&clientRequest)
	if err != nil {
		http.Error(w, "failed to parse request body", 400)
	}

	client, err := Protocol.NewClient(clientRequest.Ip, clientRequest.Port, clientRequest.Name, time.Now())
	if err != nil {
		http.Error(w, err.Error(), 400)
	}
	if client == nil {
		http.Error(w, "failed to parse request body", 400)
	}

	connectedClients[client.Id()] = client
	_, _ = fmt.Fprintf(os.Stdout, "registered client: %s with name: %s \n", client.Id(), client.Name())
	writeConnectedClients(w)
}

func updateClient(w http.ResponseWriter, r *http.Request) {
	updateRequest := updateClientRequest{}
	err := json.NewDecoder(r.Body).Decode(&updateRequest)
	if err != nil {
		http.Error(w, "failed to parse request body", 400)
	}

	client := connectedClients[updateRequest.Id]
	if client == nil {
		http.Error(w, "client are not connected", 404)
		return
	}

	client.UpdateConnected(time.Now())
	fmt.Fprintf(os.Stdout, "updated connecton to client: %s", client.Id())
	writeConnectedClients(w)
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
