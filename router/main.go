package main

import (
	"fmt"
	"net/http"
	"net/http/cgi"

	"github.com/cholthi/router/router"
	ussd "github.com/niiamon/ussd-go"
	"github.com/samora/ussd-go/sessionstores"
)

var network string = "mtn-ss"

type mtnRequest struct {
	Msisdn    string
	Message   string
	SessionId string
}

type mtnResponse struct {
	Message  string
	FreeFlow string
}

func (r mtnRequest) GetRequest() *ussd.Request {
	var ussdRequest *ussd.Request
	ussdRequest = new(ussd.Request)
	ussdRequest.Message = r.Message
	ussdRequest.Mobile = r.Msisdn

	return ussdRequest
}

func (res mtnResponse) SetResponse(r ussd.Response) {
	res.Message = r.Message
	if r.Release == true {
		res.FreeFlow = "FB"
	}
	res.FreeFlow = "FC"
}

func main() {
	err := cgi.Serve(http.HandlerFunc(handleRequest))
	if err != nil {
		fmt.Println(err)
	}
}

func handleRequest(w http.ResponseWriter, req *http.Request) {
	// mtn request adapter
	mtnR := mtnRequest{}
	mtnR.Msisdn = req.FormValue("msisdn")
	mtnR.Message = req.FormValue("INPUT")
	mtnR.SessionId = req.FormValue("sessionId")

	// mtn response adapter
	mtnRes := &mtnResponse{}

	// setup ussd object
	store := sessionstores.NewRedis("localhost:6379")
	ussd := ussd.New(store, "Router", "Dispatch")
	ussd.Ctrl(new(router.Router))

	//handle ussd Request with adapter
	ussd.Process(mtnR, mtnRes)
	// push response to http ResponseWriter
	if mtnRes.FreeFlow == "FB" {
		w.Header().Set("FreeFlow", "FB")
	} else {
		w.Header().Set("FreeFlow", "FC")
	}
	w.Write([]byte(mtnRes.Message))
}
