package main

import (
	"net/http"
	"strconv"

	"github.com/cholthi/topup/airtime"
)

func topupHandler(res http.ResponseWriter, r *http.Request) {
	reqdata, err := decodeTopupRequest(r)
	if err != nil {
		logger.Print(err)
		errres := TopUpResponse{}
		errres.StatusCode = 10
		errres.StatusMessage = "REJECTED_BUSINESS_LOGIC"
		data, _ := encodeTopupResponse(errres)
		res.Header().Set("Content-Type", "application/json")
		res.Write(data)
		return
	}

	result, err := airtime.TopUp(reqdata.Amount, reqdata.Msisdn)
	if err != nil {
		logger.Print(err)
	}

	resdata := TopUpResponse{}
	if result.ResultCode == 0 {
		fbal, err := strconv.ParseFloat(result.Balance, 32)
		if err != nil {
			logger.Print(err)
		}
		resdata.Balance = fbal
	}
	resdata.Reference = result.Reference
	resdata.StatusCode = result.ResultCode
	resdata.StatusMessage = result.Message()
	resdata.Receiver = result.SenderMsisdn

	body, _ := encodeTopupResponse(resdata)
	res.Header().Set("Content-Type", "application/json")
	res.Write(body)
}

func infoHandler(res http.ResponseWriter, r *http.Request) {
	reqdata, err := decodeInfoRequest(r)
	if err != nil {
		logger.Print(err)
		errres := InfoResponse{}
		errres.StatusCode = 10
		errres.StatusMessage = "REJECTED_BUSINESS_LOGIC"
		data, _ := encodeInfoResponse(errres)
		res.Header().Set("Content-Type", "application/json")
		res.Write(data)
		return
	}
	result, err := airtime.GetInfo(reqdata.RessellerID)
	if err != nil {
		logger.Print(err)
	}
	resdata := InfoResponse{}
	if result.ResultCode == 0 {
		/*	fbal, err := strconv.ParseFloat(result.Balance, 32)
			if err != nil {
				logger.Print(err)
			}*/
		resdata.Balance = result.Balance
	}
	resdata.Reference = result.Reference
	resdata.StatusCode = result.ResultCode
	resdata.StatusMessage = result.Message()
	resdata.Currency = result.Currency

	body, _ := encodeInfoResponse(resdata)
	res.Header().Set("Content-Type", "application/json")
	res.Write(body)
}
