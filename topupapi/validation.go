package main

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
