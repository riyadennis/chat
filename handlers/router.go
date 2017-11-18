package handlers

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Router struct {
	Router  *mux.Router
	Address string
}

func NewRouter(router *mux.Router, address string) *Router {
	return &Router{
		Router:  router,
		Address: address,
	}
}
func (r *Router) Run() {
	roomHandler := NewRoom(false)
	chatTemplateHandler := NewTemplateHandler("chat.html")
	loginTemplateHandler := NewTemplateHandler("login.html")
	loginHandler := NewLoginHandler()

	r.Router.Handle("/chat", MustAuth(chatTemplateHandler))
	r.Router.Handle("/room", roomHandler)
	r.Router.Handle("/login", loginTemplateHandler)
	r.Router.Handle("/auth/{action}/{provider}/", loginHandler)

	go roomHandler.BroadCastMessages()

	logrus.Infof("Listening to port %s", r.Address)
	err := http.ListenAndServe(r.Address, r.Router)
	if err != nil {
		logrus.Errorf("Web server run failed with error %s", err.Error())
	}
}
