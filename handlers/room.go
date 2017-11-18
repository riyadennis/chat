package handlers

import (
	"github.com/chat/entities"
	"github.com/chat/trace"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/objx"
	"net/http"
	"time"
	"os"
	"errors"
)

type Room struct {
	clients   map[*websocket.Conn]bool
	upgrader  websocket.Upgrader
	broadcast chan *entities.Message
	tracer    trace.Tracer
}

func NewRoom(tracerStatus bool) *Room {
	var tracer trace.Tracer
	if tracerStatus == true {
		tracer = trace.New(os.Stdout)
	} else {
		tracer = trace.Off()
	}
	return &Room{
		clients:   make(map[*websocket.Conn]bool),
		upgrader:  websocket.Upgrader{},
		broadcast: make(chan *entities.Message),
		tracer:    tracer,
	}
}

func (r *Room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := r.CreateClient(w, req)
	if socket != nil && err == nil {
		for {
			//@TODO need to implement readjson
			_, message, err := socket.ReadMessage()
			if err != nil {
				r.tracer.Trace("Unable to read message")
				logrus.Errorf("Unable to read message: %s", err.Error())
				break
			}

			authCookie, err := req.Cookie("auth")
			if err != nil {
				r.tracer.Trace("Invalid Login")
				logrus.Errorf("Invalid Login %s", err.Error())
				break
			}

			authCookieData := objx.MustFromBase64(authCookie.Value)
			avatarUrl := authCookieData["avatar_url"]
			jMsg := entities.NewMessage(authCookieData["name"].(string), string(message), avatarUrl.(string), time.Now())
			r.broadcast <- jMsg
		}
	}
}
func (r *Room) CreateClient(w http.ResponseWriter, req *http.Request) (*websocket.Conn, error) {
	if &r.upgrader != nil {
		socket, err := r.upgrader.Upgrade(w, req, nil)
		if err != nil {
			r.tracer.Trace("Unable to access socket")
			logrus.Errorf("Unable to access socket %s", err.Error())
			return nil, err
		}

		r.clients[socket] = true
		r.tracer.Trace("Created a new client")
		return socket,nil
	}
	return nil, errors.New("invalid request")
}
func (r *Room) BroadCastMessages() {
	for {
		message := <-r.broadcast
		r.tracer.Trace("Message:", message.GetMessage())
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
