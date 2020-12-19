package main

import "regexp"

func validateTopupRequest(data TopUpRequest) bool {
	var valid bool = true

	if data.Msisdn == "" {
		valid = false
		return valid
	}

	if data.Amount == 0 {
		valid = false
		return valid
	}

	return valid
}

func validateInfoRequest(data InfoRequest) bool {
	var valid bool = true

	if data.RessellerID == "" {
		valid = false
		return valid
	}
	return valid
}

func validateCreateAccountRequest(data CreateAccountRequest) bool {
	var valid bool = true
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if data.Email == "" {
		valid = false
	}

	if ok := emailRegex.Match([]byte(data.Email)); !ok {
		valid = false
	}

	if data.Organization == "" {
		valid = false
	}
	return valid
}

func validateCreditAccountRequest(data AccountCreditRequest) bool {
	var valid bool = true

	if data.Amount == 0.0 {
		valid = false
	}

	return valid
}
