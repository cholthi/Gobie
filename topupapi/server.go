package main

import (
	"net/http"
	"time"
)

type Server struct {
	addr    string
	handler http.Handler
}

func NewServer(addr string, handler http.Handler) *Server {
	var server *Server
	server = new(Server)
	server.addr = addr
	server.handler = handler

	return server
}

func (s Server) run() {
	server := &http.Server{
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
		Handler:      s.handler,
		Addr:         s.addr,
	}

	logger.Fatal(server.ListenAndServe())
}
