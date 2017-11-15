package entities

import "time"

type Message struct {
	Name    string
	Message string
	When    time.Time
}

func NewMessage(name, message string, when time.Time) *Message {
	return &Message{
		Name:    name,
		Message: message,
		When:    when,
	}
}
func (m *Message) GetMessage() string {
	return m.Message
}
