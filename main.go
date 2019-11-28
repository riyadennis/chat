package main

import (
	"flag"

	"github.com/riyadennis/chat/config"
	"github.com/riyadennis/chat/internal"
	"github.com/sirupsen/logrus"
)
var (
	address = flag.String("address", ":8084", "Port to which web server will listen")
	trace = flag.Bool("traceStatus", false, "error handling and tracing")
)

func main() {
	conf, err := config.ParseConfig("config.yaml")
	if err != nil {
		logrus.Errorf("invalid config :: %v", err)
	}
	err = internal.NewServer(*address, *trace, conf).Run()
	if err != nil {
		logrus.Errorf("unable to run the router :: %v", err)
	}
}
