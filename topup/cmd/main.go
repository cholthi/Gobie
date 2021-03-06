package main

import (
	"flag"
	"fmt"

	"github.com/cholthi/topup/airtime"
)

type userKey string

var UserKey userKey = "api_user"

func main() {
	var amount = flag.Float64("amount", 0, "The amount to topup the subscriber with. E.g 100. in SSP")
	var msisdn = flag.String("msisdn", "", "The phone number of the subscriber to topup")

	flag.Parse()

	//var res *airtime.TopupResponse = new(airtime.TopupResponse)

	result, err := airtime.TopUp(*amount, *msisdn)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Ref:%v\nPhone:%v\nMy Balance:%v", result.Reference, result.SenderMsisdn, result.Balance)
}
