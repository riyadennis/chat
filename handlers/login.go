package handlers

import (
	"fmt"
	"github.com/stretchr/gomniauth"
	"strings"
	"net/http"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/objx"
	"github.com/stretchr/gomniauth/common"
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
	//user, err := getUser(provider,r.URL.RawQuery)
	//authCookie := objx.New(map[string]interface{}{
	//	"name":user.Name(),
	//}).MustBase64()
	//cookie := &http.Cookie{Name: "auth", Value: authCookie,Path:"/"}

	w.Header().Set("Location", loginUrl)
	w.WriteHeader(http.StatusTemporaryRedirect)
	//http.SetCookie(w,cookie)
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
func getUser(provider string, url string) (common.User, error) {
	gProvider, err := gomniauth.Provider(provider)
	credentials, err := gProvider.CompleteAuth(objx.MustFromURLQuery(url))
	if err != nil {
		return nil, err
	}
	user, err := gProvider.GetUser(credentials)
	if err != nil {
		return nil, err
	}
	return user, nil
}
