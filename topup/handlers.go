package main

import (
	"encoding/json"
	"net/http"

	"github.com/cholthi/topupapi/model"
)

func HandleTopUp(res http.ResponseWriter, r *http.Request) {
	reqdata := topUpRequest{}
	jd := json.NewDecoder(r.Body)
	err := jd.Decode(&reqdata)

	if err != nil {
		errorResponse(res, http.StatusBadRequest, err.Error())
		return
	}

	if err = validateTopUpRequestData(reqdata); err != nil {
		errorResponse(res, http.StatusBadRequest, err.Error())
		return
	}
	apiuser := getUserFromContext(r.Context()) //panics
	if err = checkBusinessConstrains(reqdata, apiuser.ID); err != nil {
		errorResponse(res, http.StatusBadRequest, err.Error())
		return
	}
	var params map[string]interface{} = make(map[string]interface{}, 0)
	params["amount"] = reqdata.Amount
	if reqdata.Currency != "" {
		params["currency"] = reqdata.Currency
	}
	params["recipient"] = reqdata.Recipient

	err = model.AccountCommitTransaction(params)
	if err != nil {
		errorResponse(res, http.StatusInternalServerError, err.Error())
		return
	}
	// Carefull, The money is here!
	var resp Response = &Response{}
	go sendAirtime(reqdata, resp)

	res.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(res)
	enc.Encode(&resp)
	return
	// log and end request
}

func CreateAccount(res http.ResponseWriter, r *http.Request) {
	// unmarshal request
	// validate request
	// create account in db
	// prepare response (account id)
	// log and end request

}
