package handlers

import (
	"github.com/chat/trace"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Room struct {
	clients   map[*websocket.Conn]bool
	upgrader  websocket.Upgrader
	broadcast chan []byte
	tracer    trace.Tracer
}

func NewRoom() *Room {
	return &Room{
		clients:   make(map[*websocket.Conn]bool),
		upgrader:  websocket.Upgrader{},
		broadcast: make(chan []byte),
		tracer:    trace.Off(),
	}
}

func (r *Room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := r.upgrader.Upgrade(w, req, nil)
	if err != nil {
		r.tracer.Trace("Unable to access socket")
		logrus.Errorf("Unable to access socket %s", err.Error())
		return
	}
	defer socket.Close()
	r.clients[socket] = true
	r.tracer.Trace("Created a new client")
	for {
		_, message, err := socket.ReadMessage()
		if err != nil {
			r.tracer.Trace("Unable to read message")
			logrus.Errorf("Unable to read message: %s", err.Error())
			break
		}
		r.broadcast <- message
	}
}

func (r *Room) BroadCastMessages() {
	for {
		message := <-r.broadcast
		r.tracer.Trace("Message:", string(message))
		for client := range r.clients {
			client.WriteMessage(websocket.TextMessage, message)
		}
	}
}
