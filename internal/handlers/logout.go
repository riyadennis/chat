package handlers

import "net/http"

type Logout struct{}

func (l *Logout) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	cookie := &http.Cookie{
		Name:   "auth",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
	w.Header().Set("Location", "/login")
	w.WriteHeader(http.StatusTemporaryRedirect)
}
