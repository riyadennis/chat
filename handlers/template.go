package handlers

import (
	"github.com/sirupsen/logrus"
	"html/template"
	"net/http"
	"path/filepath"
	"sync"
	"github.com/stretchr/objx"
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

type Data struct {
	Host     string
	UserData string
}

func (t *TemplateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.Once.Do(func() {
		var err error
		rootPath, _ := filepath.Abs("templates")
		path := filepath.Join(rootPath, t.FileName)
		t.Template, err = template.ParseFiles(path)
		if err != nil {
			logrus.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
	if t.Template != nil {
		cookie, err := r.Cookie("auth")
		if err != nil {
			logrus.Errorf("Cant find authentication cookie returned error %s", err)
		}
		data := map[string]interface{}{
			"Host": r.Host,
		}
		if cookie != nil {
			data["UserData"] = objx.MustFromBase64(cookie.Value)
		}
		t.Template.Execute(w, data)
	}

}
