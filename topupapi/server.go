package main

import (
	"context"
	"net/http"
	"time"
)

type Server struct {
	addr    string
	handler http.Handler
	server  *http.Server
}

func NewServer(addr string, handler http.Handler) *Server {
	var server *Server
	server = new(Server)
	server.addr = addr
	server.handler = handler

	http_server := &http.Server{
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
		Handler:      handler,
		Addr:         addr,
	}

	server.server = http_server

	return server
}

func (s *Server) run() {
	logger.Fatal(s.server.ListenAndServe())
}

func (s *Server) Shutdown(ctx context.Context) {
	logger.Print(s.server.Shutdown(ctx))
}
