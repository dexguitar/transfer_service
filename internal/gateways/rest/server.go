package rest

import (
	"net/http"
)

type Server struct {
	srv *http.Server
}

func NewServer(handler http.Handler, addr string) *Server {
	return &Server{
		srv: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
	}
}

func (s *Server) Start() error {
	return s.srv.ListenAndServe()
}
