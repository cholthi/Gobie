package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
)

func decodeTopupRequest(r *http.Request) (TopUpRequest, error) {
	var decoded TopUpRequest = TopUpRequest{}
	jsondec := json.NewDecoder(r.Body)
	err := jsondec.Decode(&decoded)
	if err != nil && err != io.EOF {
		logger.Print(err)
		return decoded, err
	}
	if ok := validateTopupRequest(decoded); !ok {
		return decoded, errors.New("request data missing required fields")
	}
	return decoded, nil
}

func encodeTopupResponse(res TopUpResponse) ([]byte, error) {
	body, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func decodeInfoRequest(r *http.Request) (InfoRequest, error) {
	var decoded InfoRequest = InfoRequest{}
	jsondec := json.NewDecoder(r.Body)
	err := jsondec.Decode(&decoded)
	if err != nil && err != io.EOF {
		logger.Print(err)
		return decoded, err
	}
	if ok := validateInfoRequest(decoded); !ok {
		return decoded, errors.New("request data missing required fields")
	}
	return decoded, nil
}

func encodeInfoResponse(res InfoResponse) ([]byte, error) {
	body, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func initLogger() {
	file := "/var/log/topupapi/api.log"
	output, err := os.OpenFile(file, os.O_WRONLY|os.O_APPEND, 0666)
	//defer output.Close()
	if err != nil {
		log.Println(err)
	}

	logger.SetOutput(output)
	logger.SetPrefix("topup-api")
	logger.SetFlags(log.Lshortfile)
}

func attachMiddlewares(h http.Handler, mids ...Middleware) http.Handler {

	for _, handler := range mids {
		h = handler(h)
	}

	return h
}
