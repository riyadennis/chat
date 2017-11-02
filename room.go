package main

import (
	"net/http"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

const BufferSize = 1024

type room struct {
	//this channel holds all the messages
	forward chan []byte
	//this is the channel of clients wishing to join
	join chan *client
	//this is a channel of client wishing to leave
	leave chan *client
	//hold all active clients
	clients map[*client]bool
}

func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
}
func (r *room) Run() {
	for {
		select {
		case client := <-r.join:
			r.clients[client] = true
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.send)
		case msg := <-r.forward:
			for client := range r.clients {
				client.send <- msg
			}
		}
	}
}

var upgrader = websocket.Upgrader{ReadBufferSize: BufferSize, WriteBufferSize: BufferSize}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		logrus.Errorf("Unable to access socket %s", err.Error())
		return
	}
	client := &client{
		socket: socket,
		send:   make(chan []byte, BufferSize),
		room:   r,
	}
	r.join <- client
	defer func() { r.leave <- client }()
	client.Write()
	client.Read()
}
