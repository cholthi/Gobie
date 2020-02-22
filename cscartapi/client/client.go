package client

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"time"
)

// A client struct is an http client abstraction.
// it utilizes behind the scene the http. Client struct object
// It understand API specific semantics particular to this project.
type Client struct {
	login  string
	apiKey string
	client *http.Client
}

func (c *Client) SetLogin(login string) {
	c.login = login
}

func (c *Client) SetApiKey(key string) {
	c.apiKey = key
}

func initClient() *http.Client {
	transport := http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{Timeout: 30 * time.Second}
	client.Transport = &transport
	return client
}

func addBasicAuth(login, apikey string, r *http.Request) *http.Request {
	v := fmt.Sprintf("%s:%s", login, apikey)
	encoded := base64.StdEncoding.EncodeToString([]byte(v))

	header := "Basic " + encoded
	r.Header.Set("Authorization", header)
	return r
}

func (c *Client) Do(r *http.Request) http.Response {
	req := addBasicAuth(c.login, c.apiKey, r)
	resp, err := c.client.Do(req)
	//b, _ := ioutil.ReadAll(resp.Body)
	//fmt.Printf("%s", b)

	if err != nil {
		log.Fatalln(err)
	}
	return *resp
}

func NewClient(login, apikey string) *Client {
	client := new(Client)
	client.client = initClient()
	client.SetApiKey(apikey)
	client.SetLogin(login)

	return client
}
