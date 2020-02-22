package airtime

import (
	"bytes"
	"encoding/xml"
	"math/rand"
	"net/http"
	"time"
)

type Endpoint interface {
	PrepareRequest() ([]byte, error)
	PrepareResponse(http.Response) (interface{}, error)
	GetEndpoint() string
}

func randomString(len int) string {
	rand.Seed(time.Now().UnixNano())
	token := make([]byte, len)
	_, err := rand.Read(token)
	if err != nil {
		rand.Read(token)
	}
	return string(token)
}

func toXML(r interface{}) string {
	header := &bytes.Buffer{}
	header.Write([]byte(`<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ext="http://external.interfaces.ers.seamless.com/"><soapenv:Header/><soapenv:Body>`))
	enc := xml.NewEncoder(header)
	if err := enc.Encode(r); err != nil {
		panic(err)
	}
	header.Write([]byte(`</soapenv:Body></soapenv:Envelope>`))
	//c.log("marshalled request:", header.String())
	return header.String()
}
