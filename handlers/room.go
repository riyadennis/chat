package handlers

import (
	"github.com/chat/trace"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
	"github.com/chat/entities"
	"os"
	"time"
	"github.com/stretchr/objx"
)

type Room struct {
	clients   map[*websocket.Conn]bool
	upgrader  websocket.Upgrader
	broadcast chan *entities.Message
	tracer    trace.Tracer
}

func NewRoom() *Room {
	return &Room{
		clients:   make(map[*websocket.Conn]bool),
		upgrader:  websocket.Upgrader{},
		broadcast: make(chan *entities.Message),
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
		var jMsg entities.Message
		_, message, err := socket.ReadMessage()
		authCookie, err := req.Cookie("auth")
		if err!=nil{
			r.tracer.Trace("Invalid Login")
			logrus.Errorf("Invalid Login %s", err.Error())
			break
		}
		jMsg.Name = objx.MustFromBase64(authCookie.Value)
		jMsg.Message = string(message)
		jMsg.When = time.Now()
		if err != nil {
			r.tracer.Trace("Unable to read message")
			logrus.Errorf("Unable to read message: %s", err.Error())
			break
		}

		r.broadcast <- &jMsg
	}
}

func (r *Room) BroadCastMessages() {
	for {
		message := <-r.broadcast
		r.tracer.Trace("Message:", message.Message)
		for client := range r.clients {
			err := client.WriteJSON(message)
			if err != nil {
				r.tracer.Trace("Unable to write message")
				logrus.Errorf("Unable to write message: %s", err.Error())
				break
			}
		}
	}
}
