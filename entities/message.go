package entities

import "time"

type Message struct {
	Name    string
	Message string
	When    time.Time
	AvatarUrl string
}

func NewMessage(name, message,avatarUrl string, when time.Time) *Message {
	return &Message{
		Name:    name,
		Message: message,
		When:    when,
		AvatarUrl:avatarUrl,
	}
}
func (m *Message) GetMessage() string {
	return m.Message
}
