package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/cholthi/topupapi/model"
	"github.com/cholthi/topup/airtime"
)



func validateTopUpRequestData(r Request) error {
	var err error = nil
	mes := make([]string, 0)
	mp := map[string][]string{}

	if r.Recipient == "" {
		mes = append(mes, "Recipient field is required")
	}
	if r.Amount == 0 {
		mes = append(mes, "Amount is required or can not be zero")
	}
	if len(mes) > 1 {
		mp["validationErrors"] = mes
		data, _ := json.Marshal(mp)
		err = errors.New(string(data))
		return err
	}
	return err
}

func checkBusinessConstrains(r Request, userid int64) error {
 var err error = nil
  errs := make([]string,0)
  errormap map[string][]string
  acc, err := model.GetAccountByUserID(userid)
  if err != nil {
    return err
  }

  if acc.Balance < 0 {
    errs = append(errs, "Not enought account balance")
  }

  if acc.Balance < 3 {
    errs = append(errs, "Can not send below min amount (3 SSP)")
  }

  if acc.Balance > 50000 {
    errs = append(errs, "Can not send over max amount (50000 SSP)")
  }

  if len(errs) > 1 {
		errormap["validationErrors"] = errs
		data, _ := json.Marshal(errormap)
		err = errors.New(string(data))
		return err
	}
  return nil

}

func sendAirtime(r Request, res *Response) {
	topupres, err := airtime.TopUp(r.Amount, r.Recipient)
	if err != nil {
		panic(err)
	}

	res.Balance = topupres.Balance
	res.MSISDN = topupres.SenderMsisdn
	res.StatusCode = topupres.ResultCode
	res.Reference = topupres.Reference
	res.StatusMessage = topupres.Message()
  return
}

func getUserFromContext(ctx context.Context) model.User{
	userctx := ctx.Value(UserKey)
	u, err := userctx.(model.User)
	if err != nil {
		panic(err)
		}
		return u
}
