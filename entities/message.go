package entities

import "time"

// Message struct carries message that user broadcasts
type Message struct {
	Name      string
	Message   string
	When      time.Time
	AvatarUrl string
}

// NewMessage creates message struct from a list of params
func NewMessage(name, message, avatarUrl string, when time.Time) *Message {
	return &Message{
		Name:      name,
		Message:   message,
		When:      when,
		AvatarUrl: avatarUrl,
	}
}

// GetMessage is the getter for the message
func (m *Message) GetMessage() string {
	return m.Message
}
