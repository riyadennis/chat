package main

import (
	"net/http"
	"path/filepath"
	"github.com/sirupsen/logrus"
	"github.com/chat/lib"
	"flag"
)

func main() {
	rootPath, _ := filepath.Abs("templates")
	//making the folder templates accessible
	fs := http.FileServer(http.Dir(rootPath))
	http.Handle("/", fs)

	templateHandler := lib.NewTemplateHandler("chat.html")
	roomHandler := lib.NewRoom()
	http.Handle("/chat", lib.MustAuth(templateHandler))
	http.Handle("/room", roomHandler)

	go roomHandler.BroadCastMessages()
	var address = flag.String("address", ":8080", "Port to which webserver will listen")
	flag.Parse()
	logrus.Infof("Listening to port %s", *address)
	err := http.ListenAndServe(*address, nil)
	if err != nil {
		logrus.Errorf("Web server run failed with error %s", err.Error())
	}
}
