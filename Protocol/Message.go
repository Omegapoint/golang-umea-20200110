package Protocol

import (
	"encoding/json"
	uuid "github.com/satori/go.uuid"
	"time"
)

type Message struct {
	Id uuid.UUID
	Sent time.Time
	Message string
}

func NewMessage(id uuid.UUID, msg string) *Message {
	message := new(Message)

	message.Id = id
	message.Message = msg
	message.Sent = time.Now()

	return message
}

func (m *Message) Serialize() ([]byte, error) {
	return json.Marshal(m)
}
