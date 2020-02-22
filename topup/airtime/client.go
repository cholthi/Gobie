package airtime

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// DEFAULT_HOST server address
const DEFAULT_HOST = ""

// VERSION of this airtime client
const VERSION = "1.0"

var endpoint string
var Client *client

//Client interacts with the api server and creates response objects
type client struct {
	BaseURL    string
	Password   string
	HTTPClient http.Client
	Logger     Log
}

func init() {
	Client = NewClient("/var/log/airtime.log")
}

type Log func(...interface{})

func (c *client) log(v ...interface{}) {
	if c.Logger != nil {
		c.Logger(v...)
	}
}
func (c *client) doRequest(ctx context.Context, endpoint Endpoint) (interface{}, error) {
	//c.log("Calling method:", endpoint)
	body, err := endpoint.PrepareRequest()
	if err != nil {
		c.log("Error marshalling request", err)
	}
	// make context cancellable http call
	url := c.Endpoint(endpoint.GetEndpoint())
	bf := []byte(body)
	httpreq, err := http.NewRequest("POST", url, bytes.NewBuffer(bf))
	httpreq.Header.Set("Content-Type", "text/xml")
	if err != nil {
		c.log("Error:", err)
		return nil, err
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
		return nil, ctx.Err()
	case resp := <-ch:
		close(ch)
		//c.log("API resp for ", url, ":\n", b.String())
		// ATTENTION!
		res, err := endpoint.PrepareResponse(*resp)
		if err != nil {
			c.log("Error Unmarshalling Response", err)
		}
		//value := new(TopupResponse)
		//value.Reference = topupres.Reference
		//value.ResultCode = topupres.ResultCode
		//value.Status = topupres.Status
		//value.SenderMsisdn = topupres.SenderMsisdn
		//value.Balance = topupres.Balance

		return res, nil

	}
}
func (c *client) Endpoint(endpoint string) string {
	base := c.BaseURL
	if c.BaseURL == "" {
		base = DEFAULT_HOST
	}
	return fmt.Sprintf("%s/%s", base, endpoint)
}

func (c *client) ConfigTLS() {
	cert, err := tls.LoadX509KeyPair("/home/ubuntu/certs/cert.pem", "/home/ubuntu/certs/key.pem")
	if err != nil {
		c.log("Error:", err)
		panic(err)
	}
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}
	//tlsConfig.InsecureSkipVerify = false

	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}
	c.HTTPClient.Transport = transport

}

func TopUp(amount float64, msisdn string) (*TopupResponse, error) {

	ctx := context.Background()
	timeout := 5 * time.Second
	ctx_timeout, cancel := context.WithTimeout(ctx, timeout)

	defer cancel()
	endpoint := &TopUpEndpoint{}
	endpoint = endpoint.SetReceiver(msisdn)
	endpoint = endpoint.SetAmount(amount)
	endpoint = endpoint.SetEndpoint("topupservice/service")

	result, err := Client.doRequest(ctx_timeout, endpoint)
	if err != nil {
		Client.log("Error with sending request by client.", err)
		panic(err)
	}
	topupres, ok := result.(*TopupResponse)
	if !ok {
		Client.log("Can not assert client.DoRequest to *TopupResponse")
		panic("Can not continue. Zombies lying on the way")
	}
	return topupres, err
}

func GetInfo(resellerid string) (*InformationPrincipalResponse, error) {
	ctx := context.Background()
	timeout := time.Second * 5
	ctx_timeout, cancel := context.WithTimeout(ctx, timeout)

	defer cancel()
	endpoint := &InformationEndpoint{}
	endpoint.SetID(resellerid)
	endpoint.SetType("RESELLERID")
	endpoint.SetEdpoint("topupservice/service")

	result, err := Client.doRequest(ctx_timeout, endpoint)
	if err != nil {
		Client.log("Error with sending request by client.", err)
		panic(err)
	}
	infores, ok := result.(*InformationPrincipalResponse)
	if !ok {
		Client.log("Can not assert client.DoRequest to *InformationPrincipalResponse")
		panic("Can not continue. Zombies lying on the way")
	}
	return infores, err
}

func NewClient(file string) *client {
	var client *client = new(client)
	var out *os.File
	func() {
		o, err := os.OpenFile(file, os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			o = os.Stderr
		}
		out = o
	}()
	client.BaseURL = "https://172.16.100.97:8913"
	//client.BaseURL = "http://localhost:8913" //test server
	logger := log.New(out, "MTN-EVD", log.LstdFlags)
	client.Logger = logger.Print

	//wire http client
	client.HTTPClient = http.Client{}
	client.HTTPClient.Timeout = time.Second * 30
	client.ConfigTLS()
	return client
}
