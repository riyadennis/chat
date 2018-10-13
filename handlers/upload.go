package handlers

import (
	"net/http"
)

type uploadHandler struct{}

func (u *uploadHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	userID := req.FormValue("user_id")
	if userID == "" {
		http.Error(w, "User not found in the request", http.StatusInternalServerError)
	}
	_, _, err := req.FormFile("avatar_file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
