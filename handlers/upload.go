package handlers

import (
	"io/ioutil"
	"net/http"
	"path"
)

type uploadHandler struct{}

func (u *uploadHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	userid := req.FormValue("userid")
	file, header, err := req.FormFile("avatarfile")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	filename := path.Join("avatars", "avatar"+userid+path.Ext(header.Filename))
	err = ioutil.WriteFile(filename, data, 777)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
