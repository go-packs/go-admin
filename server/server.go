package server

import (
	"github.com/go-packs/go-admin"
	"net/http"
)

type Server struct {
	Registry *admin.Registry
	Addr     string
}

func NewServer(reg *admin.Registry, addr string) *Server {
	return &Server{Registry: reg, Addr: addr}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()
	mux.Handle("/admin/", NewRouter(s.Registry))
	
	handler := Logger(mux)
	handler = Recovery(handler)
	
	return http.ListenAndServe(s.Addr, handler)
}
