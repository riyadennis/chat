package main

import (
	"net/http"
	"path/filepath"
	"github.com/sirupsen/logrus"
	"github.com/chat/lib"
	"flag"
	"github.com/gorilla/mux"
)

func main() {
	rootPath, _ := filepath.Abs("templates")
	//making the folder templates accessible
	fs := http.FileServer(http.Dir(rootPath))
	http.Handle("/", fs)
	r := mux.NewRouter()
	chatTemplateHandler := lib.NewTemplateHandler("chat.html")
	loginTemplateHandler := lib.NewTemplateHandler("login.html")
	loginProviderHandler := lib.NewLoginProviderHandler()

	roomHandler := lib.NewRoom()
	r.Handle("/chat", lib.MustAuth(chatTemplateHandler))
	r.Handle("/login", loginTemplateHandler)
	r.Handle("/room", roomHandler)
	r.Handle("/auth/login/{provider}", loginProviderHandler)

	go roomHandler.BroadCastMessages()
	var address = flag.String("address", ":8080", "Port to which webserver will listen")
	flag.Parse()
	logrus.Infof("Listening to port %s", *address)
	err := http.ListenAndServe(*address, r)
	if err != nil {
		logrus.Errorf("Web server run failed with error %s", err.Error())
	}
}
