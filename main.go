package main

import (
	"net/http"
	"sync"
	"html/template"
	"path/filepath"
	"github.com/sirupsen/logrus"
	"github.com/chat/lib"
	"flag"
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
	t.Template.Execute(w, r)
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
	var address = flag.String("address", ":8080", "Port to which webserver will listen")
	flag.Parse()
	logrus.Infof("Listening to port %s", *address)
	err := http.ListenAndServe(*address, nil)
	if err != nil {
		logrus.Errorf("Web server run failed with error %s", err.Error())
	}
}
