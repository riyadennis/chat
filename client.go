package main

import (
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type client struct {
	socket websocket.Conn
	send   chan []byte
	room   *room
}

func (c *client) Read() {
	defer c.socket.Close()
	_, msg, err := c.socket.ReadMessage()
	if err != nil {
		logrus.Errorf("Unable to read message %s", err.Error())
	}
	room := room{}
	room.forward <- msg
}
func (c *client) Write() {
	defer c.socket.Close()
	for msg := range c.send {
		err := c.socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			logrus.Errorf("Unable to write message %s", err.Error())
		}
	}
}
