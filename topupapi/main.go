package main

import (
	"log"

	"github.com/gorilla/mux"
)

var logger log.Logger

func main() {
	var server *Server
	var middlewares []Middleware
	initLogger()
	mux := mux.NewRouter()
	mux.HandleFunc("/airtime/api/subscriber/topup", topupHandler).Methods("POST")
	mux.HandleFunc("/airtime/api/subscriber/balance", infoHandler).Methods("POST")

	middlewares = []Middleware{LoggingMiddleware(logger), AuthMiddleware(logger)}

	wrappedhandler := attachMiddlewares(mux, middlewares...)

	server = NewServer(":4445", wrappedhandler)

	//Run the Server
	server.run()
}
