package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/chat/config"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/gomniauth/common"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
)

//TODO need to add login github account
type loginHandler struct {
	Config *config.Config
}

func NewLoginHandler(config *config.Config) *loginHandler {
	return &loginHandler{
		Config: config,
	}
}

func (lh loginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	uri := strings.Split(r.URL.RequestURI(), "/")
	action := uri[2]
	provider := uri[3]
	switch action {
	case "login":
		loginUrl, err := getLoginURL(provider, lh.Config)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			logrus.Error(err)
		}
		w.Header().Set("Location", loginUrl)
		w.WriteHeader(http.StatusTemporaryRedirect)
	case "callback":
		user, err := getUser(provider, r.URL.RawQuery, lh.Config)
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
func getLoginURL(provider string, config *config.Config) (string, error) {
	loginUrl := ""
	for _, providerConf := range config.Auth.Providers {
		if provider == "google" && providerConf.Name == "google" {
			gp := google.New(providerConf.Client, providerConf.Secret, providerConf.URL)
			loginUrl, err := gp.GetBeginAuthURL(nil, nil)
			if err != nil {
				return "", err
			}
			return loginUrl, nil
		}
		if provider == "facebook" && providerConf.Name == "facebook" {
			fp := facebook.New(providerConf.Client, providerConf.Secret, providerConf.URL)
			loginUrl, err := fp.GetBeginAuthURL(nil, nil)
			if err != nil {
				return "", err
			}
			return loginUrl, nil
		}
	}

	return loginUrl, nil
}
func getUser(provider string, url string, config *config.Config) (common.User, error) {
	for _, providerConf := range config.Auth.Providers {
		if provider == "google" && providerConf.Name == "google" {
			gp := google.New(providerConf.Client, providerConf.Secret, providerConf.URL)
			credentials, err := gp.CompleteAuth(objx.MustFromURLQuery(url))
			user, err := gp.GetUser(credentials)
			if err != nil {
				return nil, err
			}
			return user, nil
		}
		if provider == "facebook" && providerConf.Name == "facebook" {
			fp := facebook.New(providerConf.Client, providerConf.Secret, providerConf.URL)
			credentials, err := fp.CompleteAuth(objx.MustFromURLQuery(url))
			user, err := fp.GetUser(credentials)
			if err != nil {
				return nil, err
			}
			return user, nil
		}
	}
	return nil, errors.New("Invalid provider")
}
func createCookieFromUser(user common.User) *http.Cookie {
	authCookie := objx.New(map[string]interface{}{
		"name":       user.Name(),
		"avatar_url": user.AvatarURL(),
	}).MustBase64()
	return &http.Cookie{Name: "auth", Value: authCookie, Path: "/"}
}
