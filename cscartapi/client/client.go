package client

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"time"
)

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
	client := &http.Client{Timeout: 30 * time.Second}
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
