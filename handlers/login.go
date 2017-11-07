package handlers

import (
	"fmt"
	"github.com/stretchr/gomniauth"
	"strings"
	"net/http"
)

type loginHandler struct{}

func NewLoginHandler() *loginHandler {
	return &loginHandler{}
}

func (lh loginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	uri := strings.Split(r.URL.RequestURI(), "/")
	action := uri[2]
	provider := uri[3]
	if action != "login" {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Authentication action %s is not supported", action)
	}
	gprovider, err := gomniauth.Provider(provider)
	if err != nil {
		http.Error(w, fmt.Sprintf("Access the provider %s encountered error %s", provider, err.Error()), http.StatusInternalServerError)
		return
	}
	loginUrl, err := gprovider.GetBeginAuthURL(nil, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error when trying to begin URL %s", loginUrl), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Location", loginUrl)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

