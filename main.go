package main

import (
	"flag"
	"github.com/chat/config"
	"github.com/chat/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"path/filepath"
)

func main() {
	rootPath, _ := filepath.Abs("templates")
	//making the folder templates accessible
	fs := http.FileServer(http.Dir(rootPath))
	http.Handle("/templates", fs)

	r := mux.NewRouter()
	chatTemplateHandler := handlers.NewTemplateHandler("chat.html")
	loginTemplateHandler := handlers.NewTemplateHandler("login.html")
	loginProviderHandler := handlers.NewLoginProviderHandler()

	roomHandler := handlers.NewRoom()
	r.Handle("/chat", handlers.MustAuth(chatTemplateHandler))
	r.Handle("/login", loginTemplateHandler)
	r.Handle("/room", roomHandler)
	r.Handle("/auth/login/{provider}", loginProviderHandler)

	go roomHandler.BroadCastMessages()
	var address = flag.String("address", ":8080", "Port to which web server will listen")
	flag.Parse()
	logrus.Infof("Listening to port %s", *address)
	conf := config.ParseConfig("config.yaml")
	handlers.SetupAuth(conf)
	err := http.ListenAndServe(*address, r)
	if err != nil {
		logrus.Errorf("Web server run failed with error %s", err.Error())
	}
}
