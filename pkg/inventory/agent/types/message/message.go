package message

import (
	"encoding/json"

	"github.com/stmcginnis/gofish"
)

type Message struct {
	Endpoint string          `json:"endpoint,omitempty"`
	Session  *gofish.Session `json:"session,omitempty"`

	Additional string `json:"additional,omitempty"`
}

func Parse(data []byte) (m *Message, err error) {
	return m, json.Unmarshal(data, &m)
}

func (m *Message) Marshal() ([]byte, error) {
	return json.Marshal(m)
}
