package app

import (
	"github.com/julienschmidt/httprouter"
	"library-go/pkg/logging"
	"net/http"
)

type Server struct {
	Logger *logging.Logger
	config *Config
	router *httprouter.Router
}

func (s *Server) New(config *Config, router *httprouter.Router, logger *logging.Logger) {
	s.config = config
	s.router = router
	s.Logger = logger
}

func (s *Server) Start() error {
	s.Logger.Infof("Server starts at %v", s.config.ServerAddress())
	if err := http.ListenAndServe(s.config.ServerAddress(), s.router); err != nil {
		s.Logger.Fatalf("Server start fail. err %v", err)
	}
	return nil
}
