package lib

import (
	"github.com/gorilla/websocket"
	"net/http"
	"github.com/sirupsen/logrus"
	"github.com/chat/trace"
	"os"
)

type Room struct {
	Clients   map[*websocket.Conn]bool
	Upgrader  websocket.Upgrader
	Broadcast chan []byte
	Tracer trace.Tracer
}

func NewRoom() *Room {
	return &Room{
		Clients:   make(map[*websocket.Conn]bool),
		Upgrader:  websocket.Upgrader{},
		Broadcast: make(chan []byte),
		Tracer: trace.New(os.Stdout),
	}
}

func (r *Room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := r.Upgrader.Upgrade(w, req, nil)
	if err != nil {
		r.Tracer.Trace("Unable to access socket")
		logrus.Errorf("Unable to access socket %s", err.Error())
		return
	}
	defer socket.Close()
	r.Clients[socket] = true
	r.Tracer.Trace("Created a new client")
	for {
		_, message, err := socket.ReadMessage()
		if err != nil {
			r.Tracer.Trace("Unable to read message")
			logrus.Errorf("Unable to read message: %s", err.Error())
			break
		}
		r.Broadcast <- message
	}
}

func (r *Room) BroadCastMessages() {
	for {
		message := <-r.Broadcast
		r.Tracer.Trace("Message:",string(message))
		for client := range r.Clients {
			client.WriteMessage(websocket.TextMessage, message)
		}
	}
}
