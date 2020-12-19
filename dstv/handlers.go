package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/cholthi/dstv/mailc"
)

const apikey = "test"

type APIRespond struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

type RechargeRequest struct {
	CardNumber string `json:"smartCardNumber"`
	TxnId      string `json:"transactionId"`
	Amount     string `json:"amount"`
	PhoneNo    string `json:"phoneNumber"`
}

func getPackage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		query := r.URL.Query()
		serial := query.Get("serial")
		seriali, err := strconv.Atoi(serial)
		if err != nil {
			errRespond(w, err.Error())
		}
		pack := getPackageBySerial(seriali)
		body := marshal(pack)
		_, err = w.Write(body)
		if err != nil {
			errRespond(w, err.Error())
			logger.Println(err)
		}
		return
	}
}

func getPackageBySerial(i int) Package {
	packages := getPackages()
	for _, p := range packages {
		if p.SerialNo == i {
			return p
		}
	}
	return packages[0]
}
func apiPackages(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		data := getPackages()
		body := marshal(data)
		_, err := w.Write(body)
		if err != nil {
			errRespond(w, "Internal server error")
			return
		}
		return
	}
	errRespond(w, "Invalid request method")
	return
}

//Recharge handler writes req to storage and notify backend office

func apiRecharge(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		pack, err := unMarshalBody(r.Body)
		if err != nil {
			errRespond(w, "Invalid request body params")
			logger.Println(err)
			return
		}
		err = store.Persist(pack)
		if err != nil {
			errRespond(w, "Internal server error")
			logger.Println(err)
			return
		}
		//notify
		go func() {
			attachment := csvdatabase
			mailer := mailc.NewSMTP(logger, SMTP_HOST, SMTP_PASSWORD, SMTP_USERNAME, SMTP_PORT)
			addr := []string{"swangin@dstvsouthsudan.com", "awak@dstvsouthsudan.com", "kevin@dstvsouthsudan.com", "anyang@targetmedia.biz"}
			body := makeBody(pack)
			err := mailer.Send("New Mgurush Customer Payment", body, attachment, addr)
			if err != nil {
				logger.Panicln(err)
			}
		}()

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"error":false,"message":"recharge OK"}`))
		return
	}
	errRespond(w, "Invalid request Method")
	return
}

func authenticateMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("x-api-key")
		if key == "" || key != apikey {
			errRespond(w, "Invalid API key")
			return
		}
		h.ServeHTTP(w, r)
	})
}

func marshal(apires interface{}) []byte {
	data, err := json.Marshal(apires)
	if err != nil {
		logger.Panicln(err)
	}
	return data
}

func unMarshalBody(r io.Reader) (RechargeRequest, error) {
	dec := json.NewDecoder(r)
	recha := RechargeRequest{}
	err := dec.Decode(&recha)

	if err != nil {
		return RechargeRequest{}, err
	}
	return recha, nil
}
func errRespond(w http.ResponseWriter, message string) {
	errres := APIRespond{true, message}
	w.Header().Set("Content-Type", "application/json")
	w.Write(marshal(errres))
}

func makeBody(p RechargeRequest) string {
	//buf := strings.Builder{}
	sc := fmt.Sprintf("Customer Smart Card number %s\n Phone Number %s\n Amount Paid %s\n Mgurush TransactionID %s", p.CardNumber, p.PhoneNo, p.Amount, p.TxnId)
	return sc
}
