package main

import (
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	server := setup()
	server.ListenAndServe()
}

func getAddr() string {
	if port := os.Getenv("PORT"); port != "" {
		return ":" + port
	}

	return ":8913"
}

func setup() *http.Server {
	return &http.Server{
		Addr:         getAddr(),
		Handler:      getMux(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}

func getMux() http.Handler {
	mux := mux.NewRouter()
	mux.HandleFunc("/topupservice/service", topup)
	mux.HandleFunc("/topupservice/service2", info)
	return mux
}

func topup(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "text/xml")
	data := `<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
<soap:Body>
<ns2:requestTopupResponse xmlns:ns2="http://external.interfaces.ers.seamless.com/">
<return>
<ersReference>2015061114441812901004879</ersReference>
<resultCode>0</resultCode>
<resultDescription>SUCCESS</resultDescription>
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
