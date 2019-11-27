package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"github.com/riyadennis/chat/config"
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
		pr := checkProvider(provider, lh.Config.Auth.Providers)
		if pr == nil{
			http.Error(w, "invalid provider", http.StatusBadRequest)
		}
		l, err := loginUrl(pr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			logrus.Error(err)
		}
		w.Header().Set("Location", l)
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

func loginUrl(pr *config.Provider) (string, error){
	switch pr.Name {
	case "google": gp := google.New(pr.Client, pr.Secret, pr.URL)
		loginUrl, err := gp.GetBeginAuthURL(nil, nil)
		if err != nil {
			return "", err
		}
		return loginUrl, err
	case "facebook":
		fp := facebook.New(pr.Client, pr.Secret, pr.URL)
		loginUrl, err := fp.GetBeginAuthURL(nil, nil)
		if err != nil {
			return "", err
		}
		return loginUrl, err
	}
	return "", fmt.Errorf("cannot fetch login url for %v", pr)
}

func checkProvider(provider string, providers []*config.Provider ) *config.Provider{
	for _, pr := range providers {
		if pr.Name == provider{
			return pr
		}
	}
	return nil
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
