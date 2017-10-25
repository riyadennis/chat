package main

import (
	"net/http"
	"sync"
	"html/template"
	"os"
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
		pwd, _ := os.Getwd()
		t.Template = template.Must(template.ParseFiles(pwd + "/chatapp/templates/" + t.FileName))
	})
	t.Template.Execute(w, nil)
}

func main() {
	templateHandler := NewTemplateHandler("chat.html")
	http.Handle("/", templateHandler)
	http.ListenAndServe(":8080", nil)
}
