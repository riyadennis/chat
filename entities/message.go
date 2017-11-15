package entities

import "time"

type Message struct {
	Name    string
	Message string `json:Message`
	When    time.Time
}

func (m *Message) GetMessage() string{
	return m.Message
}
