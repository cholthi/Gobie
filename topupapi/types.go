package main

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
	Reference     string  `json:"reference"`
	StatusCode    int     `json:"statusCode"`
	StatusMessage string  `json:"statusMessage"`
	Currency      string  `json:"currency,omitempty"`
	Balance       float64 `json:"balance,omitempty"`
}
