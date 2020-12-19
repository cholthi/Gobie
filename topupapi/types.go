package main

import (
	"github.com/gobuffalo/uuid"
)

type TopUpRequest struct {
	Msisdn   string  `json:"msisdn"`
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency,omitempty"`
}

type TopUpResponse struct {
	Reference     string  `json:"reference"`
	StatusCode    int     `json:"statusCode"`
	StatusMessage string  `json:"statusMessage"`
	Receiver      string  `json:"receiver,omitempty"`
	Balance       float64 `json:"balance,omitempty"`
}

type InfoRequest struct {
	RessellerID string `json:"resellerid"`
	Currency    string `json:"currency"`
}

type InfoResponse struct {
	//Reference     string  `json:"reference"`
	StatusCode    int     `json:"statusCode"`
	StatusMessage string  `json:"statusMessage"`
	Currency      string  `json:"currency,omitempty"`
	Balance       float64 `json:"balance,omitempty"`
}

type CreateAccountRequest struct {
	Organization string `json:"organization"`
	Email        string `json:"email"`
}

type AccountCreateResponse struct {
	Message   string    `json:"message"`
	AccountID uuid.UUID `json:"account_id"`
}

type AccountCreditRequest struct {
	Organization string  `json:"organization"`
	Amount       float64 `json:"amount"`
	Currency     string  `json:"currency,omitempty"`
}

type AccountCreditResponse struct {
	Message string  `json:"message"`
	Balance float64 `json:"balance"`
}
