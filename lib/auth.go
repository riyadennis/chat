package lib

import (
	"net/http"
	"strings"
	"fmt"
	"github.com/stretchr/gomniauth"
	"github.com/chat/config"
	"github.com/stretchr/gomniauth/providers/google"
)

type authHandler struct {
	next http.Handler
}
type loginProviderHandler struct{}

func NewLoginProviderHandler() *loginProviderHandler {
	return &loginProviderHandler{}
}
func (ah *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("auth")
	httpStatus := CheckCookie(w, r)
	if httpStatus != 0 {
		http.Error(w, err.Error(), httpStatus)
		return
	}
	ah.next.ServeHTTP(w, r)
}
func CheckCookie(w http.ResponseWriter, r *http.Request) int {
	_, err := r.Cookie("auth")
	if err == http.ErrNoCookie {
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
		return http.StatusTemporaryRedirect
	}
	return 0
}
func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}
func SetupAuth(conf *config.Config) {
	gomniauth.SetSecurityKey(conf.Auth.Security)
	for _, provider := range conf.Auth.Providers {
		fmt.Print(provider.Name)
		if provider.Name == "google" {
			gomniauth.WithProviders(
				google.New(provider.Client, provider.Secret, provider.Url),
			)
		}
	}
}
func (lh loginProviderHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
