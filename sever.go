package library_app

import (
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(port string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}
	log.Info("Start Server")
	return s.httpServer.ListenAndServe()
}
