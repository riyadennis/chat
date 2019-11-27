package internal

import (
	"github.com/gorilla/mux"
	"github.com/riyadennis/chat/config"
	"github.com/riyadennis/chat/handlers"
	"net/http"
	"path/filepath"
)

type Server struct {
	Address string
	Trace bool
}

func (s *Server) Run() error{
	conf, err := config.ParseConfig("config.yaml")
	if err != nil {
		return err
	}
	r := mux.NewRouter()
	router := handlers.NewRouter(r, s.Address, conf)
	router.Run(s.Trace)
}

func LoadTemplates(path string) error{
	rootPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	//making the folder templates accessible
	fs := http.FileServer(http.Dir(rootPath))
	http.Handle("/templates", fs)
	return nil
}