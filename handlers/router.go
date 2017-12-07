package handlers

import (
	"net/http"

	"github.com/chat/config"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type Router struct {
	Router  *mux.Router
	Address string
	Config  *config.Config
}

func NewRouter(router *mux.Router, address string, config *config.Config) *Router {
	return &Router{
		Router:  router,
		Address: address,
		Config:  config,
	}
}

func (r *Router) Run(tracerStatus bool) {
	roomHandler := NewRoom(tracerStatus)
	chatTemplateHandler := NewTemplateHandler("chat.html")
	loginTemplateHandler := NewTemplateHandler("login.html")
	loginHandler := NewLoginHandler(r.Config)
	logoutHandler := &Logout{}

	r.Router.Handle("/chat", MustAuth(chatTemplateHandler))
	r.Router.Handle("/room", roomHandler)
	r.Router.Handle("/login", loginTemplateHandler)
	r.Router.Handle("/logout", logoutHandler)
	r.Router.Handle("/auth/{action}/{provider}/", loginHandler)

	go roomHandler.BroadCastMessages()

	logrus.Infof("Listening to port %s", r.Address)
	err := http.ListenAndServe(r.Address, r.Router)
	if err != nil {
		logrus.Errorf("Web server run failed with error %s", err.Error())
	}
}
