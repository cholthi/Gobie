package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cholthi/topup/airtime"
)

func TestTopup(t *testing.T) {

	result, err := airtime.TopUp(1, "211925415377")
	if err != nil {
		t.Error(err)
	}
	if result.Reference != "2015061114441812901004879" {
		t.Errorf("Wanted %s, Got %s", "2015061114441812901004879", result.Reference)
	}

}

func TestInfo(t *testing.T) {

	result, err := airtime.GetInfo("211925415377")

	if err != nil {
		t.Error(err)
	}
	if result.Reference != "2015061114325676001004877" {
		t.Errorf("Wanted %s, Got %s", "2015061114325676001004877", result.Reference)
	}

	if result.Balance != 586.00 {
		t.Errorf("Wanted %f, Got %f", float64(586.00), result.Balance)
	}

}

func TestTopupHandler(t *testing.T) {
	var body *bytes.Buffer = bytes.NewBuffer([]byte(`{"msisdn":"0925415377","amount":10}`))
	req, err := http.NewRequest("POST", "/airtime/api/subscriber/topup", body)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(topupHandler)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Wanted http code %d, Got http code %d", http.StatusOK, rr.Code)
	}

	res := decodetopupResponse(ioutil.NopCloser(rr.Body))
	if res.Reference != "2015061114441812901004879" {
		t.Errorf("Wanted %s, Got %s", "2015061114441812901004879", res.Reference)
	}

}

func TestInfoHandler(t *testing.T) {
	var body *bytes.Buffer = bytes.NewBuffer([]byte(`{"resellerid":"211925415377","currency":"SSP"}`))
	req, err := http.NewRequest("POST", "/airtime/api/subscriber/balance", body)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(infoHandler)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Wanted http code %d, Got http code %d", http.StatusOK, rr.Code)
	}

	res := decodeinfoResponse(ioutil.NopCloser(rr.Body))
	if res.Reference != "2015061114325676001004877" {
		t.Errorf("Wanted %s, Got %s", "2015061114325676001004877", res.Reference)
	}
	//t.Error(res)
}

func decodetopupResponse(r io.ReadCloser) TopUpResponse {
	var res TopUpResponse = TopUpResponse{}
	dec := json.NewDecoder(r)
	dec.Decode(&res)
	return res

}

func decodeinfoResponse(r io.ReadCloser) InfoResponse {
	var res InfoResponse = InfoResponse{}
	dec := json.NewDecoder(r)
	dec.Decode(&res)
	return res

}
