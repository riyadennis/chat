package main

import (
	"context"
	"flag"

	"github.com/riyadennis/chat/config"
	"github.com/riyadennis/chat/handlers"
	"github.com/riyadennis/chat/internal"
	"github.com/sirupsen/logrus"
)

var ctx context.Context

func main() {
	address := flag.String("address", ":8080", "Port to which web server will listen")
	trace := flag.Bool("traceStatus", false, "error handling and tracing")
	flag.Parse()

	conf, err := config.ParseConfig("config.yaml")
	if err != nil {
		logrus.Errorf("invalid config :: %v", err)
	}
	err = internal.LoadTemplates(conf.TemplatePath)
	if err != nil {
		logrus.Errorf("unable to open template files :: %v", err)
	}

	handlers.SetupAuth(conf.Auth)
	s := &internal.Server{Address: *address, Trace: *trace}
	err = s.Run()
	if err != nil {
		logrus.Errorf("unable to run the router :: %v", err)
	}
}
