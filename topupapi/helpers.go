package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cholthi/topupapi/model"
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

func decodeCreditAccountRequest(r *http.Request) (AccountCreditRequest, error) {
	var decoded AccountCreditRequest = AccountCreditRequest{}
	jsondec := json.NewDecoder(r.Body)
	err := jsondec.Decode(&decoded)
	if err != nil {
		logger.Print(err)
		return decoded, err
	}
	if ok := validateCreditAccountRequest(decoded); !ok {
		return decoded, errors.New("request data missing required fields")
	}
	return decoded, nil
}

func encodeCreditAccountResponse(res AccountCreditResponse) ([]byte, error) {
	body, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func errorResponse(res http.ResponseWriter, code int, errorvalue interface{}) {
	genericerror := struct {
		Error   bool
		Message string
	}{}
	if err, ok := errorvalue.(error); ok {
		genericerror.Error = true
		genericerror.Message = err.Error()
	} else {
		genericerror.Error = true
		genericerror.Message = errorvalue.(string)
	}
	body, _ := json.Marshal(genericerror)
	res.WriteHeader(code)
	res.Write(body)

}

func encodeInfoResponse(res InfoResponse) ([]byte, error) {
	body, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func decodeCreateAccountRequest(r *http.Request) (CreateAccountRequest, error) {
	var decoded CreateAccountRequest = CreateAccountRequest{}
	jsondec := json.NewDecoder(r.Body)
	err := jsondec.Decode(&decoded)
	if err != nil && err != io.EOF {
		logger.Print(err)
		return decoded, err
	}
	if ok := validateCreateAccountRequest(decoded); !ok {
		return decoded, errors.New("request data missing required fields")
	}
	return decoded, nil
}

func encodeCreateAccountResponse(res AccountCreateResponse) ([]byte, error) {
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

func processFailedTransaction() {
	//retrytx := []RetryTransaction{}
	retrys, err := model.GetTransactionsForRetry() //culprit
	if err != nil {
		panic(err)
	}

	for _, tnx := range retrys {
		var errs chan error = make(chan error)
		go retryTransaction(tnx, errs)
		err := <-errs
		logger.Println(err)
	}
}

func waitForTermination(done <-chan struct{}) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	select {
	case sig := <-signals:
		logger.Printf("Triggering shutdown from signal %s", sig)
	case <-done:
		logger.Print("Shutting down...")
	}
}

func scheduler(done chan struct{}) {
	ticker := time.NewTicker(time.Hour * 1)

	for {
		select {
		case t := <-ticker.C:
			logger.Printf("Retrying Failed Transaction at %s\n", t)
			processFailedTransaction()
		case <-done:
			ticker.Stop()
			logger.Println("Exiting Scheduler")
			break
		}
	}
}
