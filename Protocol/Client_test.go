package Protocol

import (
	"testing"
	"time"
)

const ip = "192.168.1.1"
const port = 8080
const name = "Some name"

func TestCreateAValidClient(t *testing.T) {
	_, err := newClient(ip, port, name, time.Now())
	if err != nil {
		t.Error("failed to instantiate a valid client")
	}
}

func TestCanCreateAClientWithPort(t *testing.T) {
	_, err := newClient(ip, port, name, time.Now())
	if err != nil {
		t.Error("failed to instantiate a valid client")
	}
}

func TestIpWithPortIsNotValid(t *testing.T) {
	_, err := newClient(ip + ":8080", port, name, time.Now())
	if err == nil {
		t.Error("port in ip string is not valid")
	}
}

func TestCorrectlyMarshalsCLient(t *testing.T) {
	created, err := time.Parse(time.RFC3339, "2019-12-19T20:10:22.525293+01:00")
	if err != nil {
		t.Error("failed to parse time")
	}
	c, err := newClient(ip, port, name, created)
	if err != nil {
		t.Error("failed to create client")
	}

	expected := "{\"ip\":\"192.168.1.1\",\"port\":8080,\"connected\":\"2019-12-19T20:10:22.525293+01:00\",\"name\":\"Some name\"}"
	actual_bytes, err := c.MarshalJSON()
	if err != nil {
		t.Error("failed to marshal json")
	}

	actual := string(actual_bytes[:])
	if expected != actual {
		t.Errorf("expected: %v, got: %v", expected, actual)
	}
}
