package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var logger log.Logger

const (
	REJECTED_AMOUNT   = 11
	REJECETED_PAYMENT = 12
	SUCCESS           = 0
)

func main() {
	done := make(chan struct{})
	defer close(done)

	//scheduler retries failed transactions with error REJECTED_AMOUNT status every 1 Hour
	go scheduler(done)
	var server *Server
	var middlewares []Middleware
	initLogger()
	authhandler, err := getAuthHandler()
	if err != nil {
		panic(err)
	}
	router := mux.NewRouter()
	//requireAuth := mux.MiddlewareFunc(requireAthenticationMiddleware(logger))
	//router.Use(requireAuth)
	router.Handle("/airtime/api/subscriber/topup", requireAthenticationMiddleware(logger)(http.HandlerFunc(topupHandler))).Methods("POST")
	router.Handle("/airtime/api/subscriber/balance", requireAthenticationMiddleware(logger)(http.HandlerFunc(infoHandler))).Methods("POST")
	router.Handle("/airtime/api/subscriber/balancev2", requireAthenticationMiddleware(logger)(http.HandlerFunc(checkBalanceHandler))).Methods("GET")
	router.Handle("/airtime/api/subscriber/account", requireAthenticationMiddleware(logger)(http.HandlerFunc(createAccountHandler))).Methods("POST")
	router.Handle("/airtime/api/subscriber/account/buycredits", requireAthenticationMiddleware(logger)(http.HandlerFunc(buyAccountCredits))).Methods("POST")
	router.PathPrefix("/auth").Handler(http.StripPrefix("/auth", authhandler))
	middlewares = []Middleware{LoggingMiddleware(logger)}

	wrappedhandler := attachMiddlewares(router, middlewares...)

	server = NewServer(":4443", wrappedhandler)

	go func() {
		waitForTermination(done)
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		server.Shutdown(ctx)
	}()

	//Run the Server
	server.run()
}
