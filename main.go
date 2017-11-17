package main

import (
	"flag"
	"github.com/chat/config"
	"github.com/chat/handlers"
	"github.com/gorilla/mux"
	"net/http"
	"path/filepath"
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
	conf := config.ParseConfig("config.yaml")
	handlers.SetupAuth(conf)
	router.Run()
}
