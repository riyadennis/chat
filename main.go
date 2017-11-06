package main

import (
	"net/http"
	"sync"
	"html/template"
	"path/filepath"
	"github.com/sirupsen/logrus"
	"github.com/gorilla/websocket"
	"fmt"
	"github.com/chat/lib"
)

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
	roomHandler := lib.NewRoom()
	http.Handle("/chat", templateHandler)
	http.Handle("/room", roomHandler)

	go roomHandler.BroadCastMessages()
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		logrus.Errorf("Web server run failed with error %s", err.Error())
	}
}
