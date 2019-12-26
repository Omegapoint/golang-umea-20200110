package Protocol

import (
	"encoding/json"
	"errors"
	uuid "github.com/satori/go.uuid"
	"net"
	"time"
)

type Client struct {
	id uuid.UUID
	ip net.IP
	port uint16
	connected time.Time
	name string
}

func NewClient(ip string, port uint16, name string, created time.Time) (*Client, error) {
	parsedIp := net.ParseIP(ip)
	if parsedIp == nil {
		return  nil, errors.New("invalid IP provided")
	}

	err := validateName(name)
	if err != nil {
		return nil, err
	}

	c := new(Client)

	id, _ := uuid.NewV4()
	c.id = id
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

func (c *Client) Port() uint16 {
	return c.port
}

func (c *Client) Connected() time.Time {
	return c.connected
}

func (c *Client) Name() string {
	return c.name
}

func (c *Client) Id() uuid.UUID {
	return c.id
}

func (c *Client) UpdateConnected(t time.Time) {
	c.connected = t
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