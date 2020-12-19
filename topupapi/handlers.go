package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/cholthi/topup/airtime"
	"github.com/cholthi/topupapi/model"
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
	u := getUser(r.Context())
	logger.Printf("%v", u)
	account, err := model.FindAccountByUserID(u.ID)
	if err != nil {
		errorResponse(res, http.StatusInternalServerError, err)
		return
	}
	if account.Balance >= reqdata.Amount {
		result, err := airtime.TopUp(reqdata.Amount, reqdata.Msisdn)
		if err != nil {
			logger.Print(err)
			//check error is REJECTED_PAYMENT or REJECTED_AMOUNT or REJECTED_TOPUP
			//requeue the top up for retry later on
		}

		tnx, err := model.NewTransaction(reqdata.Msisdn, reqdata.Amount, account.ID, result.Reference)
		if err != nil {
			errorResponse(res, http.StatusInternalServerError, err)
			return
		}
		if result.Code() == REJECTED_AMOUNT || result.Code() == REJECETED_PAYMENT {
			_ = tnx.SetStatus(model.TNX_PENDING)
			rt, err := model.NewRetryTransaction(tnx)
			logger.Printf("queued Transaction # %s for retry\n", tnx.ID)
			if err == nil {
				rt.IncrementRetry()
				err = rt.Save()
				if err != nil {
					logger.Printf("Error queuing transaction for retry: %v", err.Error())
				}
			}
		}
		err = account.DoTransaction(tnx)
		if err != nil {
			errorResponse(res, http.StatusInternalServerError, err)
			return
		}
		resdata := TopUpResponse{}
		// we consider REJECTED_AMOUNT and REJECTED_PAYMENT statuses successeful transactions because the transactions will be retried in thebackground
		if result.ResultCode == SUCCESS || result.ResultCode == REJECTED_AMOUNT || result.ResultCode == REJECETED_PAYMENT {
			resdata.Balance = account.Balance
		}
		resdata.Reference = result.Reference
		resdata.StatusCode = result.ResultCode
		if result.ResultCode == REJECTED_AMOUNT || result.ResultCode == REJECETED_PAYMENT {
			resdata.StatusMessage = "SUCCESS"
		} else {
			resdata.StatusMessage = result.Message()
		}
		resdata.Receiver = result.SenderMsisdn

		body, _ := encodeTopupResponse(resdata)
		res.Header().Set("Content-Type", "application/json")
		res.Write(body)
	} else {
		errorResponse(res, http.StatusForbidden, "Account balance amount less than topup amount ")
	}
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
	//resdata.Reference = result.Reference
	resdata.StatusCode = result.ResultCode
	resdata.StatusMessage = result.Message()
	resdata.Currency = result.Currency

	body, _ := encodeInfoResponse(resdata)
	res.Header().Set("Content-Type", "application/json")
	res.Write(body)
}

func checkBalanceHandler(res http.ResponseWriter, r *http.Request) {
	u := getUser(r.Context())
	account, err := model.FindAccountByUserID(u.ID)
	if err != nil {
		errorResponse(res, http.StatusNotFound, fmt.Errorf("%w Error: No account found", err))
		return
	}
	resdata := InfoResponse{}
	resdata.StatusCode = 0
	resdata.Balance = account.Balance
	resdata.Currency = "SSP"

	body, _ := encodeInfoResponse(resdata)
	res.Header().Set("Content-Type", "application/json")
	res.Write(body)
}

func createAccountHandler(res http.ResponseWriter, r *http.Request) {
	data, err := decodeCreateAccountRequest(r)
	u := getUser(r.Context())
	if err != nil {
		errorResponse(res, http.StatusBadRequest, err)
		return
	}

	if valid := validateCreateAccountRequest(data); !valid {
		errorResponse(res, http.StatusBadRequest, errors.New("The Email or organization can't be empty"))
		return
	}
	//check the user is admin creating accounts
	if !u.IsSuperAdmin {
		errorResponse(res, http.StatusUnauthorized, errors.New("This operation is not permitted"))
		return
	}

	account, err := model.CreateAccount(data.Email, data.Organization)
	if err != nil {
		errorResponse(res, http.StatusInternalServerError, err)
		return
	}
	if err := account.Save(); err != nil {
		errorResponse(res, http.StatusInternalServerError, err)
	}
	var resdata AccountCreateResponse = AccountCreateResponse{}
	resdata.AccountID = account.ID
	resdata.Message = "Account created successefully"

	jsonres, err := encodeCreateAccountResponse(resdata)
	if err != nil {
		errorResponse(res, http.StatusInternalServerError, err)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.Write(jsonres)

}

func retryTransaction(rt model.RetryTransaction, errs chan error) {
	tnx := rt.Transaction
	result, err := airtime.TopUp(tnx.Amount, tnx.Recipient)
	if err != nil {
		errs <- fmt.Errorf("%w error from airtime.Topup function", err)
		goto done
	}
	if result.Code() == 0 {
		account := &model.Account{}
		err := model.DB.Find(account, tnx.AccountID)
		if err != nil {
			errs <- fmt.Errorf("%w problem getting acccount from db", err)
		}
		if err := tnx.SetStatus(model.TNX_SUCCESS); err != nil {
			errs <- fmt.Errorf("%w error setting transaction status", err)
			goto done
		}
		account.SetBalance(-tnx.Amount)
		err = account.Update("amount")
		if err != nil {
			errs <- fmt.Errorf("%w error updating account in db", err)
			goto done
		}
		err = model.DB.Destroy(rt)
		if err != nil {
			errs <- fmt.Errorf("%w error removing retry transaction in db", err)
			goto done
		}
	}
done:
}

func buyAccountCredits(res http.ResponseWriter, r *http.Request) {
	reqdata, err := decodeCreditAccountRequest(r)
	if err != nil {
		logger.Print(err)
		errres := AccountCreditResponse{}
		errres.Message = "REJECTED_BUSINESS_LOGIC"
		data, _ := encodeCreditAccountResponse(errres)
		res.Header().Set("Content-Type", "application/json")
		res.Write(data)
		return
	}
	if ok := validateCreditAccountRequest(reqdata); !ok {
		errorResponse(res, http.StatusBadRequest, errors.New("Invalid request to for credit account"))
		return
	}
	u := getUser(r.Context())
	if !u.IsSuperAdmin {
		errorResponse(res, http.StatusUnauthorized, errors.New("This operation is not permitted"))
		return
	}

	account, err := model.FindAccountByOrganization(reqdata.Organization)
	if err != nil {
		errorResponse(res, http.StatusInternalServerError, err)
		return
	}
	newbalance := account.Balance + reqdata.Amount
	account.SetBalance(newbalance)
	err = account.Update("balance")
	if err != nil {
		errorResponse(res, http.StatusInternalServerError, err)
		return
	}

	resdata := AccountCreditResponse{}
	resdata.Balance = account.Balance
	resdata.Message = "Account credit recharge successfully"
	jsonres, err := encodeCreditAccountResponse(resdata)
	if err != nil {
		errorResponse(res, http.StatusInternalServerError, err)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.Write(jsonres)
}
