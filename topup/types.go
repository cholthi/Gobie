package main

type topUpRequest struct {
	Recipient string  `json:"recipient"`
	Amount    float64 `json:"amount"`
	Currency  string  `json:"currency"`
}

type topUpResponse struct {
	Reference     string  `json:"reference"`
	StatusCode    int     `json:"status_code"`
	StatusMessage string  `json:"status_message"`
	MSISDN        string  `json:"receive_by"`
	Balance       float64 `json:"account_balance"`
}

type createAccountRequest struct {
	Organisation string `json:"organisation"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Confirm      string `json:"password_comfirm"`
}
