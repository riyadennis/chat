package handlers

import (
	"github.com/gorilla/mux"
	"net/http"
	"github.com/sirupsen/logrus"
)

type Router struct {
	Router *mux.Router
	Address string
}

func (r *Router) Run(){
	roomHandler := NewRoom()
	chatTemplateHandler := NewTemplateHandler("chat.html")
	loginTemplateHandler := NewTemplateHandler("login.html")
	loginHandler := NewLoginHandler()

	r.Router.Handle("/chat", MustAuth(chatTemplateHandler))
	r.Router.Handle("/room", roomHandler)
	r.Router.Handle("/login", loginTemplateHandler)
	r.Router.Handle("/auth/{action}/{provider}/", loginHandler)

	go roomHandler.BroadCastMessages()

	err := http.ListenAndServe(r.Address, r.Router)
	if err != nil {
		logrus.Errorf("Web server run failed with error %s", err.Error())
	}
}
