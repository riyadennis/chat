package entities

import "time"

type Message struct {
	Name    map[string]interface{}
	Message string
	When    time.Time
}

func (m *Message) GetMessage() string{
	return m.Message
}
