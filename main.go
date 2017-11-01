package main

import (
	"net/http"
	"sync"
	"html/template"
	"fmt"
	"path/filepath"
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
	templateHandler := NewTemplateHandler("chat.html")
	http.Handle("/", templateHandler)
	err :=http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}
