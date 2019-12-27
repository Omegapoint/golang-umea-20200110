package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Omegapoint/golang-umea-20200110/Protocol"
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

var clients map[uuid.UUID]Protocol.Client
func main() {
	conf := getConfig()
	fmt.Printf("Using config: %v\n", conf)
	subscribe(conf)
}

func subscribe(conf config) {
	ip, err := getLocalIp()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to retreive local ip: %v", err)
		os.Exit(1)
	}

	reqBody, err := json.Marshal(map[string]interface{} {
		"ip": ip,
		"port": conf.Port,
		"name": conf.Name,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create request: %v", err)
		os.Exit(1)
	}

	nameServerUrl := "http://" + conf.NameServerIp.String() + ":" + strconv.FormatUint(uint64(conf.NameServerPort), 10)
	registerUrl := nameServerUrl + "/client"

	resp, err := http.Post(registerUrl,
		"application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		print(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		print(err)
	}
	fmt.Println(string(body))

	//go updateSubscription()
}

func updateSubscription() {

}
