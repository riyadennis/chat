package lib

import (
	"github.com/gorilla/websocket"
	"net/http"
	"github.com/sirupsen/logrus"
)

type Room struct {
	Clients   map[*websocket.Conn]bool
	Upgrader  websocket.Upgrader
	Broadcast chan []byte
}

func NewRoom() *Room {
	return &Room{
		Clients:   make(map[*websocket.Conn]bool),
		Upgrader:  websocket.Upgrader{},
		Broadcast: make(chan []byte),
	}
}

func (r *Room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := r.Upgrader.Upgrade(w, req, nil)
	if err != nil {
		logrus.Errorf("Unable to access socket %s", err.Error())
		return
	}
	defer socket.Close()
	r.Clients[socket] = true
	for {
		_, message, err := socket.ReadMessage()
		if err != nil {
			logrus.Errorf("Unable to read message: %s", err.Error())
			break
		}
		r.Broadcast <- message
	}
}

func (r *Room) BroadCastMessages() {
	for {
		message := <-r.Broadcast
		for client := range r.Clients {
			client.WriteMessage(websocket.TextMessage, message)
		}
	}
}
