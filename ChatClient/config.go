package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type config struct {
	Port uint16 `json:"port"`
	Name string `json:"name"`
	NameServerPort uint16 `json:"nameServerPort"`
	NameServerIp net.IP `json:"nameServerIp"`
}

func getConfig() config {
	return readConfigFile()
}

func getLocalIp() (net.IP, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP, nil
			}
		}
	}

	return nil, errors.New("no local ip found")
}

func readConfigFile() config {
	filePath, _ := filepath.Abs(filepath.Dir(""))
	jsonFile, err := os.Open(filePath + "/ChatClient/config.json")
	if err != nil {
		return readConfigFromStdin()
	}
	defer jsonFile.Close()

	bytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return readConfigFromStdin()
	}

	var conf config
	json.Unmarshal(bytes, &conf)
	return conf
}

func readConfigFromStdin() config {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("What is your name: ")
	name, _ := reader.ReadString('\n')

	fmt.Print("What port to use: ")
	port, _ := reader.ReadString('\n')

	numericPort, err := strconv.ParseUint(strings.TrimSpace(port), 10, 16)
	if err != nil {
		fmt.Println("Invalid port")
		return readConfigFromStdin()
	}

	fmt.Print("Ip of name server: ")
	ipString, err := reader.ReadString('\n')
	if err != nil {
		return readConfigFromStdin()
	}

	ip := net.ParseIP(strings.TrimSpace(ipString))
	if ip == nil {
		fmt.Println("Invalid ip")
		return readConfigFromStdin()
	}

	fmt.Print("Name server port: ")
	port, _ = reader.ReadString('\n')
	numericNameServerPort, err := strconv.ParseUint(strings.TrimSpace(port), 10, 16)

	return config{
		Port: uint16(numericPort),
		Name: strings.TrimSpace(name),
		NameServerIp: ip,
		NameServerPort: uint16(numericNameServerPort),
	}
}
