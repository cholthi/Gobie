package airtime

import (
	"bytes"
	"encoding/xml"
	"io"
	"net/http"
	"time"
)

type TopUpEndpoint struct {
	Req         TopUpRequest
	Res         TopupResponse
	receiver    string
	topUpAmount float64
	endpoint    string
}

func (tue *TopUpEndpoint) SetAmount(amount float64) *TopUpEndpoint {
	tue.topUpAmount = amount
	return tue
}

func (tue *TopUpEndpoint) SetReceiver(receiver string) *TopUpEndpoint {
	tue.receiver = receiver
	return tue
}

func (tue *TopUpEndpoint) SetEndpoint(endpoint string) *TopUpEndpoint {
	tue.endpoint = endpoint
	return tue
}

func (info *TopUpEndpoint) GetEndpoint() string {
	return info.endpoint
}

func (tue *TopUpEndpoint) PrepareRequest() ([]byte, error) {
	timeout := 5 * time.Microsecond
	user := User{ID: "", Type: "RESELLERUSER", RequestType: "DMark"}
	//user := User{ID: "211925415377", Type: "RESELLERUSER", RequestType: "DMark"}
	user.XMLName = xml.Name{Local: "initiatorPrincipalId"}
	ref := randomString(8)
	context := Context{"WEBSERVICE", "DMARK", timeout, user, "", ref, "dmark topup subscriber account"}
	account := new(Account)
	account.XMLName = xml.Name{Local: "senderAccountSpecifier"}
	account.ID = ""
	account.Type = "RESELLER"
	subscriber := User{XMLName: xml.Name{Local: "topupPrincipalId"}, ID: tue.receiver, Type: "SUBSCRIBERMSISDN"}
	sub_account := new(Account)
	sub_account.XMLName = xml.Name{Local: "topupAccountSpecifier"}
	sub_account.ID = tue.receiver //important! airtime receiver is taken here
	sub_account.Type = "AIRTIME"
	amt := Amount{XMLName: xml.Name{Local: "amount"}, Value: tue.topUpAmount, Currency: "SSP"}

	var req *TopUpRequest = new(TopUpRequest)
	req.Context = context
	req.Amount = amt
	req.ProductID = "TOPUP_A"
	//req.Sender = user
	req.SenderAccount = *account
	req.Subscriber = subscriber
	req.SubscriberAccount = *sub_account

	body := toXML(req)

	return []byte(body), nil

}

func (tue *TopUpEndpoint) PrepareResponse(r http.Response) (interface{}, error) {
	var resp TopupResponse = TopupResponse{}
	//footer := "</soap:Body></soap:Envelope>"
	//lastindex := strings.Index(b, footer)
	//b = b[80:lastindex] // Hard coded xml header length. How did I find this? never mind

	var b bytes.Buffer
	if _, err := io.Copy(&b, r.Body); err != nil {
		return nil, err
	}

	if err := xml.Unmarshal([]byte(b.String()), &resp); err != nil {
		//c.log("Error:", err)
		return nil, err
	}
	if resp.ResultCode != 0 {
		var _ ErrorResponse = (*TopupResponse)(nil)
		return &resp, nil
	}
	return &resp, nil
}
