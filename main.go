package main

import (
	"net/http"
	"sync"
	"html/template"
	"path/filepath"
	"github.com/sirupsen/logrus"
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
	room := newRoom()
	templateHandler := NewTemplateHandler("chat.html")
	http.Handle("/", templateHandler)
	http.Handle("/room", room)
	go room.Run()
	err :=http.ListenAndServe(":8080", nil)
	if err != nil {
		logrus.Errorf("Web server run failed with error %s", err.Error())
	}
}
