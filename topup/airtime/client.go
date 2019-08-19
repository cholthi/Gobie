package airtime

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

// DEFAULT_HOST server address
const DEFAULT_HOST = ""

// VERSION of this airtime client
const VERSION = "1.0"

var endpoint string

//Client interacts with the api server and creates response objects
type Client struct {
	BaseURL    string
	Password   string
	HTTPClient http.Client
	Logger     Log
}

type Log func(...interface{})

func (c *Client) log(v ...interface{}) {
	if c.Logger != nil {
		c.Logger(v...)
	}
}
func (c *Client) doRequest(ctx context.Context, res interface{}, options map[string]interface{}) error {
	msisdn := options["msisdn"]
	amount := options["amount"]
	endpoint := options["endpoint"]
	c.log("Calling method:", endpoint)
	req := c.prepareReq(msisdn.(string), amount.(float64))
	body := c.toXML(*req)

	// make context cancellable http call
	url := c.Endpoint(endpoint.(string))
	bf := []byte(body)
	httpreq, err := http.NewRequest("POST", url, bytes.NewBuffer(bf))
	if err != nil {
		c.log("Error:", err)
		return err
	}
	httpreq = httpreq.WithContext(ctx)

	//send our request off
	ch := make(chan *http.Response, 1)

	go func() {
		resp, err := c.HTTPClient.Do(httpreq)
		if err != nil {
			c.log("error:", err)
			ch <- nil
		}
		ch <- resp
	}()

	select {
	case <-ctx.Done():
		<-ch
		return ctx.Err()
	case resp := <-ch:
		close(ch)
		var b bytes.Buffer
		if _, err := io.Copy(&b, resp.Body); err != nil {
			return err
		}
		c.log("API resp for ", url, ":", b)
		// ATTENTION!
		topupres := c.prepareResp(b.String())
		value, _ := res.(*TopupResponse)
		value.Reference = topupres.Reference
		value.ResultCode = topupres.ResultCode
		value.Status = topupres.Status
		value.SenderMsisdn = topupres.SenderMsisdn
		value.Balance = topupres.Balance

		if err := resp.Body.Close(); err != nil {
			return err
		}

		return nil

	}
}
func (c *Client) Endpoint(endpoint string) string {
	base := c.BaseURL
	if c.BaseURL == "" {
		base = DEFAULT_HOST
	}
	return fmt.Sprintf("%s/%s", base, endpoint)
}

func (c *Client) toXML(r interface{}) string {
	header := &bytes.Buffer{}
	header.Write([]byte(`<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/"xmlns:ext="http://external.interfaces.ers.seamless.com/"><soapenv:Header/><soapenv:Body>`))
	enc := xml.NewEncoder(header)
	if err := enc.Encode(r); err != nil {
		c.log("Error:", err)
		panic(err)
	}
	header.Write([]byte(`</soapenv:Body></soapenv:Envelope>`))
	c.log("marshalled request:", header.String())
	return header.String()
}

func (c *Client) prepareReq(msisdn string, amount float64) *TopUpRequest {
	timeout := 5 * time.Microsecond
	user := User{ID: "RES0000059747", Type: "RESELLERUSER", RequestType: "DMark"}
	user.XMLName = xml.Name{Local: "senderPrincipalId"}
	ref := randomString(8)
	context := Context{"webservice", "DMARK", timeout, user, "DM@rk321", ref, "dmark", "dmark topup subscriber account"}
	account := new(Account)
	account.XMLName = xml.Name{Local: "senderAccountSpecifier"}
	account.ID = "211925415377"
	account.Type = "RESELLER"
	subscriber := User{XMLName: xml.Name{Local: "topupPrincipalId"}, ID: msisdn, Type: "SUBSCRIBERMSISDN"}
	sub_account := new(Account)
	sub_account.XMLName = xml.Name{Local: "topupAccountSpecifier"}
	sub_account.ID = msisdn
	sub_account.Type = "AIRTIME"
	amt := Amount{XMLName: xml.Name{Local: "amount"}, Value: amount}

	var req *TopUpRequest = new(TopUpRequest)
	req.Context = context
	req.Amount = amt
	req.ProductID = "TOPUP_A"
	req.Sender = user
	req.SenderAccount = *account
	req.Subscriber = subscriber
	req.SubscriberAccount = *sub_account

	return req

}

func (c *Client) prepareResp(b string) *TopupResponse {
	var resp TopupResponse
	b = b[83:] // Hard coded xml header length. How did I find this? never mind

	if err := xml.Unmarshal([]byte(b), &resp); err != nil {
		c.log("Error:", err)
		return nil
	}
	if resp.ResultCode != 0 {
		var _ ErrorResponse = (*TopupResponse)(nil)
		return &resp
	}
	return &resp
}

func randomString(len int) string {
	randomInt := func(min, max int) int {
		return min + rand.Intn(max-min)
	}
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(randomInt(65, 90))
	}
	return string(bytes)
}

func (c *Client) ConfigTLS() {
	cert, err := tls.LoadX509KeyPair("/home/cholthi/certs/cert.pem", "/home/cholthi/certs/key.pem")
	if err != nil {
		c.log("Error:", err)
		panic(err)
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	transport := &http.Transport{TLSClientConfig: tlsConfig}
	c.HTTPClient.Transport = transport

}

func TopUp(amount float64, msisdn string, result interface{}) error {
	logger := log.New(os.Stdout, "Airtime-API", log.Lshortfile)
	var client *Client = new(Client)
	client.ConfigTLS()
	params := map[string]interface{}{
		"amount": amount,
		"msisdn": msisdn,
	}
	params["endpoint"] = "topupservice/service"
	client.BaseURL = "https://172.16.100.97:8913"
	client.Logger = logger.Print
	ctx := context.Background()
	timeout := 3 * time.Second
	ctx_timeout, cancel := context.WithTimeout(ctx, timeout)

	defer cancel()

	err := client.doRequest(ctx_timeout, result, params)
	return err
}
