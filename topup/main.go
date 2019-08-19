package main

import (
	"flag"
	"fmt"

	"github.com/cholthi/topup/airtime"
)

func main() {
	var amount = flag.Float64("amount", 0, "The amount to topup the subscriber with. E.g 100. in SSP")
	var msisdn = flag.String("msisdn", "21175415377", "The phone number of the subscriber to topup")

	flag.Parse()

	var res *airtime.TopupResponse = new(airtime.TopupResponse)

	err := airtime.TopUp(*amount, *msisdn, res)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%+v", res)
}
