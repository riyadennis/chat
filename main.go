package main

import (
	"net/http"
	"sync"
	"html/template"
	"path/filepath"
	"github.com/sirupsen/logrus"
	"github.com/gorilla/websocket"
	"fmt"
)

var clients = make(map[*websocket.Conn]bool)
var broadCasts = make(chan []byte)
var upgrader = websocket.Upgrader{}

type TemplateHandler struct {
	Once     sync.Once
	FileName string
	Template *template.Template
}

func NewTemplateHandler(fileName string) *TemplateHandler {
	return &TemplateHandler{
		FileName: fileName,
	}
}

func (t *TemplateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.Once.Do(func() {
		rootPath, _ := filepath.Abs("templates")
		path := filepath.Join(rootPath, t.FileName)
		t.Template = template.Must(template.ParseFiles(path))
	})
	t.Template.Execute(w, nil)
}

func main() {
	rootPath, _ := filepath.Abs("templates")
	//making the folder templates accessible
	fs := http.FileServer(http.Dir(rootPath))
	http.Handle("/", fs)

	templateHandler := NewTemplateHandler("chat.html")
	http.Handle("/chat", templateHandler)

	http.HandleFunc("/room", HandleConnection)
	go ReadMessages()
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		logrus.Errorf("Web server run failed with error %s", err.Error())
	}
}

func HandleConnection(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		logrus.Errorf("Unable to access socket %s", err.Error())
		return
	}
	defer socket.Close()
	clients[socket] = true
	for {
		_, message, err := socket.ReadMessage()
		if err != nil {
			logrus.Errorf("Unable to read message: %s", err.Error())
			break
		}
		broadCasts <- message
	}
}

func ReadMessages() {
	for {
		message := <-broadCasts
		fmt.Println(string(message))
		for client := range clients {
			client.WriteMessage(websocket.TextMessage, message)
		}
	}
}
