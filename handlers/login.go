package handlers

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/common"
	"github.com/stretchr/objx"
	"net/http"
	"strings"
)

//TODO need to add login github account
type loginHandler struct{}

func NewLoginHandler() *loginHandler {
	return &loginHandler{}
}

func (lh loginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	uri := strings.Split(r.URL.RequestURI(), "/")
	action := uri[2]
	provider := uri[3]
	switch action {
	case "login":
		loginUrl, err := getLoginURL(provider)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			logrus.Error(err)
		}
		w.Header().Set("Location", loginUrl)
		w.WriteHeader(http.StatusTemporaryRedirect)
	case "callback":
		user, err := getUser(provider, r.URL.RawQuery)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			logrus.Error(err)
		}
		if user != nil {
			cookie := createCookieFromUser(user)
			http.SetCookie(w, cookie)
			w.Header().Set("Location", "/chat")
			w.WriteHeader(http.StatusTemporaryRedirect)
		}

	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Auth action %s not supported", action)
	}

}

/**
Gets the login url from the provider
*/
func getLoginURL(provider string) (string, error) {
	gProvider, err := gomniauth.Provider(provider)
	if err != nil {
		return "", err
	}
	loginUrl, err := gProvider.GetBeginAuthURL(nil, nil)
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
func createCookieFromUser(user common.User) *http.Cookie {
	
	authCookie := objx.New(map[string]interface{}{
		"name": user.Name(),
		"avatar_url": user.AvatarURL(),
	}).MustBase64()
	return &http.Cookie{Name: "auth", Value: authCookie, Path: "/"}
}
