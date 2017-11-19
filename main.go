package main

import (
	"flag"
	"github.com/chat/config"
	"github.com/chat/handlers"
	"github.com/gorilla/mux"
	"net/http"
	"path/filepath"
	"github.com/sirupsen/logrus"
)

func main() {
	rootPath, _ := filepath.Abs("templates")
	//making the folder templates accessible
	fs := http.FileServer(http.Dir(rootPath))
	http.Handle("/templates", fs)

	r := mux.NewRouter()
	var address = flag.String("address", ":8080", "Port to which web server will listen")
	flag.Parse()
	router := handlers.NewRouter(r, *address)

	conf, err := config.ParseConfig("config.yaml")
	if err!=nil{
		logrus.Errorf("Invalid config %s", err.Error())
	}
	handlers.SetupAuth(conf)

	router.Run()
}
