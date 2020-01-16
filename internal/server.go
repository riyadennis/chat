package internal

import (
	"github.com/gorilla/mux"
	"github.com/riyadennis/chat/config"
	"github.com/riyadennis/chat/internal/handlers"
	"github.com/sirupsen/logrus"
	"net/http"
	"path/filepath"
)

type Server struct {
	Address string
	Trace bool
	Config *config.Config
}

func NewServer(address string, trace bool, conf *config.Config) *Server{
	return &Server{
		Address: address,
		Trace:   trace,
		Config:  conf,
	}
}

func (s *Server) Run() error{
	err := LoadTemplates(s.Config.TemplatePath)
	if err != nil {
		logrus.Errorf("unable to open template files :: %v", err)
	}
	handlers.SetupAuth(s.Config.Auth)
	r := mux.NewRouter()
	router := handlers.NewRouter(r, s.Address, s.Config)
	router.Run(s.Trace)
	return nil
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