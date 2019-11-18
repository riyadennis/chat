package main

import (
	"context"
	"flag"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/riyadennis/chat/config"
	"github.com/riyadennis/chat/handlers"
	"github.com/sirupsen/logrus"
)

var ctx context.Context

func main() {
	rootPath, err := filepath.Abs("templates")
	if err != nil {
		logrus.Errorf("unable to open template files")
	}
	//making the folder templates accessible
	fs := http.FileServer(http.Dir(rootPath))
	http.Handle("/templates", fs)

	conf, err := config.ParseConfig("config.yaml")
	if err != nil {
		logrus.Errorf("invalid config %s", err.Error())
	}

	r := mux.NewRouter()
	var address = flag.String("address", ":8080", "Port to which web server will listen")
	router := handlers.NewRouter(r, *address, conf)

	handlers.SetupAuth(conf)

	trace := flag.Bool("traceStatus", false, "error handling and tracing")
	flag.Parse()

	router.Run(*trace)
}
