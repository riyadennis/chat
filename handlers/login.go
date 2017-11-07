package handlers

import (
	"fmt"
	"github.com/stretchr/gomniauth"
	"strings"
	"net/http"
	"github.com/sirupsen/logrus"
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
	loginUrl, err := getLoginURL(provider)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logrus.Error(err)
	}
	w.Header().Set("Location", loginUrl)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

/**
Gets the login url from the provider
 */
func getLoginURL(provider string) (string, error) {
	gprovider, err := gomniauth.Provider(provider)
	if err != nil {
		return "", err
	}
	loginUrl, err := gprovider.GetBeginAuthURL(nil, nil)
	if err != nil {
		return "", err
	}
	return loginUrl, nil
}
