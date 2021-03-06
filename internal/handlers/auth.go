package handlers

import (
	"github.com/riyadennis/chat/internal/config"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/common"
	"net/http"
)

type authHandler struct {
	next http.Handler
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

// CheckCookie checks whether auth cookie is already set or not
func CheckCookie(w http.ResponseWriter, r *http.Request) int {
	_, err := r.Cookie("auth")
	if err == http.ErrNoCookie {
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
		return http.StatusTemporaryRedirect
	}
	return 0
}

// MustAuth function calls real handler after checking authentication
func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

// SetupAuth will check security key and then create providers
func SetupAuth(auth *config.Auth) {
	gomniauth.SetSecurityKey(auth.Security)
	var providerList []common.Provider
	for _, provider := range auth.Providers {
		providerList = append(providerList, provider.GetProvider())
	}
	gomniauth.WithProviders(providerList...)

}
