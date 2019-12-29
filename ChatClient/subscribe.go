package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

func subscribe(conf config) chan ClientMap {
	ip, err := getLocalIp()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to retreive local ip: %v\n", err)
		os.Exit(1)
	}

	reqBody, err := json.Marshal(map[string]interface{} {
		"ip": ip,
		"port": conf.Port,
		"name": conf.Name,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create request: %v\n", err)
		os.Exit(1)
	}

	nameServerUrl := "http://" + conf.NameServerIp.String() + ":" + strconv.FormatUint(uint64(conf.NameServerPort), 10)
	registerUrl := nameServerUrl + "/client"

	resp, err := http.Post(registerUrl,
		"application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to register client: %v\n", err)
		os.Exit(1)
	}

	defer resp.Body.Close()
	clientId := parseIdResponse(resp)

	subscription := make(chan ClientMap)
	go updateSubscription(registerUrl, subscription, clientId)
	return subscription
}

func parseResponse(resp *http.Response) ClientMap {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read response: %v\n", err)
		os.Exit(1)
	}

	var clients ClientMap
	_ = json.Unmarshal(body, &clients)
	return clients
}

func parseIdResponse(resp *http.Response) uuid.UUID {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read response: %v\n", err)
		os.Exit(1)
	}

	var idResponse struct {
		Id string `json:"id"`
	}
	_ = json.Unmarshal(body, &idResponse)
	id, _ := uuid.FromString(idResponse.Id)
	return id
}

func updateSubscription(nameServerUrl string, channel chan ClientMap, id uuid.UUID) {
	for true {
		channel <- makeUpdateRequest(nameServerUrl, id)
		time.Sleep(time.Minute)
	}
}

func makeUpdateRequest(nameServerUrl string, id uuid.UUID) ClientMap {
	reqBody := []byte(fmt.Sprintf(`{"id": "%s"}`, id.String()))
	req, _ := http.NewRequest(http.MethodPatch, nameServerUrl, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "lost connection to name server: %v\n", err)
		return nil
	}
	defer resp.Body.Close()

	return parseResponse(resp)
}
