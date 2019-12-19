package main

import (
	"encoding/json"
	"errors"
	"net"
	"time"
)

type Client struct {
	ip net.IP
	port uint16
	connected time.Time
	name string
}

func newClient(ip string, port uint16, name string, created time.Time) (*Client, error) {
	parsedIp := net.ParseIP(ip)
	if parsedIp == nil {
		return  nil, errors.New("invalid IP provided")
	}

	err := validateName(name)
	if err != nil {
		return nil, err
	}

	c := new(Client)
	c.ip = parsedIp
	c.port = port
	c.connected = created
	c.name = name

	return c, nil
}

func validateName(name string) error {
	if len(name) == 0 {
		return errors.New("name is required")
	}

	if len(name) > 140 {
		return errors.New("name can not be longer than 140 characters")
	}

	return nil
}

func (c *Client) Ip() net.IP {
	return c.ip
}

func (c *Client) Connected() time.Time {
	return c.connected
}

func (c *Client) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Ip net.IP `json:"ip"`
		Port uint16 `json:"port"`
		Connected time.Time `json:"connected"`
		Name string `json:"name"`
	}{
		Ip: c.ip,
		Port: c.port,
		Connected: c.connected,
		Name: c.name,
	})
}