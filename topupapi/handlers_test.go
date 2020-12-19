package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cholthi/topup/airtime"
)

func TestTopup(t *testing.T) {

	serverTest := httptest.NewServer(http.HandlerFunc(topup))
	defer serverTest.Close()

	airtime.Client = airtime.NewClient("./airtime.log")
	airtime.Client.BaseURL = serverTest.URL
	result, err := airtime.TopUp(1, "211925415377")
	if err != nil {
		t.Error(err)
	}
	if result.Reference != "2015061114441812901004879" {
		t.Errorf("Wanted %s, Got %s", "2015061114441812901004879", result.Reference)
	}

}

func TestInfo(t *testing.T) {

	serverTest := httptest.NewServer(http.HandlerFunc(info))
	defer serverTest.Close()

	airtime.Client = airtime.NewClient("./airtime.log")
	airtime.Client.BaseURL = serverTest.URL
	result, err := airtime.GetInfo("211925415377")

	if err != nil {
		t.Error(err)
	}
	if result.Reference != "2015061114325676001004877" {
		t.Errorf("Wanted %s, Got %s", "2015061114325676001004877", result.Reference)
	}

	if result.Balance != 586.00 {
		t.Errorf("Wanted %f, Got %f", float64(586.00), result.Balance)
	}

}

/*func TestTopupHandler(t *testing.T) {
	var body *bytes.Buffer = bytes.NewBuffer([]byte(`{"msisdn":"0925415377","amount":10}`))
	req, err := http.NewRequest("POST", "/airtime/api/subscriber/topup", body)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(topupHandler)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Wanted http code %d, Got http code %d", http.StatusOK, rr.Code)
	}

	res := decodetopupResponse(ioutil.NopCloser(rr.Body))
	if res.Reference != "2015061114441812901004879" {
		t.Errorf("Wanted %s, Got %s", "2015061114441812901004879", res.Reference)
	}

}

func TestInfoHandler(t *testing.T) {
	var body *bytes.Buffer = bytes.NewBuffer([]byte(`{"resellerid":"211925415377","currency":"SSP"}`))
	req, err := http.NewRequest("POST", "/airtime/api/subscriber/balance", body)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(infoHandler)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Wanted http code %d, Got http code %d", http.StatusOK, rr.Code)
	}

	res := decodeinfoResponse(ioutil.NopCloser(rr.Body))
	if res.Reference != "2015061114325676001004877" {
		t.Errorf("Wanted %s, Got %s", "2015061114325676001004877", res.Reference)
	}
	t.Error(res)
}*/
func topup(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "text/xml")
	data := `<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
<soap:Body>
<ns2:requestTopupResponse xmlns:ns2="http://external.interfaces.ers.seamless.com/">
<return>
<ersReference>2015061114441812901004879</ersReference>
<resultCode>11</resultCode>
<resultDescription>REJECTED_AMOUNT</resultDescription>
<requestedTopupAmount>
<currency>SSP</currency>
<value>1</value>
</requestedTopupAmount>
<senderPrincipal>
<principalId>
<id>211966220019</id>
<type>RESELLERID</type>
</principalId>
<principalName>Test Reseller</principalName>
<accounts>
<account>
<accountDescription>Reseller account for reseller</accountDescription>
<accountSpecifier>
<accountId>211966220019</accountId>
<accountTypeId>RESELLER</accountTypeId>
</accountSpecifier>
<balance>
<currency>SSP</currency>
<value>491.00000</value>
</balance>
<creditLimit>
<currency>SSP</currency>
<value>0.00000</value>
</creditLimit>
</account>
</accounts>
<status>Active</status>
<msisdn>211966220019</msisdn>
</senderPrincipal>
<topupAccountSpecifier>
<accountId>211965615584</accountId>
<accountTypeId>AIRTIME</accountTypeId>
</topupAccountSpecifier>
<topupAmount>
<currency>SSP</currency>
<value>1.00</value>
</topupAmount>
<topupPrincipal>
<principalId>
<id>211965615584</id>
<type>SUBSCRIBERID</type>
</principalId>
<principalName/>
<accounts>
<account>
<accountSpecifier>
<accountId>211965615584</accountId>
<accountTypeId>AIRTIME</accountTypeId>
</accountSpecifier>
<balance>
<currency>SSP</currency>
<value>223.07</value>
</balance>
<creditLimit>
<currency>SSP</currency>
<value>0</value>
</creditLimit>
</account>
</accounts>
</topupPrincipal>
</return>
</ns2:requestTopupResponse>
</soap:Body>
</soap:Envelope>`

	rw.Write([]byte(data))
}

func info(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "text/xml")
	data := `<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
<soap:Body>
<ns2:requestPrincipalInformationResponse xmlns:ns2="http://external.interfaces.ers.seamless.com/">
<return>
  <ersReference>2015061114325676001004877</ersReference>
<resultCode>0</resultCode>
<resultDescription>SUCCESS</resultDescription>
<requestedPrincipal>
<principalId>
<id>211966220019</id>
<type>RESELLERID</type>
</principalId>
<principalName>Test Reseller</principalName>
<accounts>
<account>
<accountDescription>RESELLER</accountDescription>
<accountSpecifier>
<accountId>211966220019</accountId>
<accountTypeId>RESELLER</accountTypeId>
</accountSpecifier>
<balance>
<currency>SSP</currency>
<value>586.00000</value>
</balance>
<creditLimit>
<currency>SSP</currency>
<value>0.00000</value>
</creditLimit>
</account>
</accounts>
<status>Active</status>
<msisdn>211966220019</msisdn>
</requestedPrincipal>
</return>
</ns2:requestPrincipalInformationResponse>
</soap:Body>
</soap:Envelope>`

	rw.Write([]byte(data))
}

func decodetopupResponse(r io.ReadCloser) TopUpResponse {
	var res TopUpResponse = TopUpResponse{}
	dec := json.NewDecoder(r)
	dec.Decode(&res)
	return res

}

func decodeinfoResponse(r io.ReadCloser) InfoResponse {
	var res InfoResponse = InfoResponse{}
	dec := json.NewDecoder(r)
	dec.Decode(&res)
	return res

}
