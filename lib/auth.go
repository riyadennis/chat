package lib

import "net/http"

type authHandler struct {
	next http.Handler
}
type Auth struct {
	name string
	value string
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
func MustAuth(handler http.Handler) http.Handler{
	return &authHandler{next:handler}
}
