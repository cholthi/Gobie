package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"

	c "github.com/cholthi/cscartapi/client"
)

// This package deals with CScart API specific interactions.

const BASE_URL = "https://agoro.co/api.php"

var errorCscartResponse error = errors.New("Error from API:")
var errorInvalidMethod error = errors.New("Error: Invalid Method")

var client *c.Client = c.NewClient("chol@dmarkmobile.com", "r476310Y887705A4648V2o7BZ3R818cP")

type CscartResponse struct {
	Data interface{}
}

func Api(method, resource, body string, query url.Values) *CscartResponse {
	if method == "" {
		log.Fatalln(errorInvalidMethod)
	}
	//fmt.Println(body)
	buf := bytes.NewBuffer([]byte(body))
	url := getUrl()
	query.Set("_d", resource)

	req, err := http.NewRequest(method, url, buf)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		log.Fatal(err)
	}
	q := req.URL.Query()
	for k, v := range query {
		for _, value := range v {
			q.Add(k, value)
		}
	}
	req.URL.RawQuery = q.Encode()
	resp := client.Do(req)
	defer resp.Body.Close()
	csresp := parseBody(resp.Body)
	return csresp
}

func getUrl() string {
	return BASE_URL
}

func parseBody(r io.Reader) *CscartResponse {

	var apiresponse interface{} //Attention
	decoder := json.NewDecoder(r)
	decoder.Decode(&apiresponse)
	ret := new(CscartResponse)
	ret.Data = apiresponse
	return ret
}
