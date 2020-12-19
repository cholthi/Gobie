package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/cholthi/cscartapi/api"
)

type Work struct {
	Product     string              `json:"product"`
	Description string              `json:"full_description"`
	Amount      string              `json:"price"`
	Imagepairs  []string            `json:"image_pairs"`
	Options     map[string][]string `json:"options,omitempty"`
}

var category string
var vendor string
var file string

var exchangeRate float64 = 36
var margin float64 = 0.5

var wg sync.WaitGroup = sync.WaitGroup{}

func main() {
	jobs := make(chan Work, 1000)
	results := make(chan api.CscartResponse, 1000)

	flag.StringVar(&category, "category", "", "Provide category to add products to on cscart")
	flag.StringVar(&vendor, "vendor", "1", "Company id to add products under")
	flag.StringVar(&file, "file", "", "inventory file containing products to upload")
	flag.Float64Var(&exchangeRate, "rate", 1, "the Rate to any currency to UGX")
	flag.Float64Var(&margin, "margin", 0.3, "margin as a percentage of products price")
	flag.Parse()

	var products []Work = make([]Work, 0)
	fileobject, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	jsonstr, err := ioutil.ReadAll(fileobject)
	if err != nil {
		log.Fatal(err)
	}
	for i := 1; i < 10; i++ {
		wg.Add(1)
		go worker(i, jobs, results)
	}
	fmt.Println("decoding json file")
	err = json.Unmarshal(jsonstr, &products)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for _, work := range products {
		jobs <- work
	}
	close(jobs)
	for result := range results {
		fmt.Printf("%v", result)
	}
	fmt.Println()
	fmt.Println("Am Done!")

}

func worker(id int, jobs <-chan Work, results chan<- api.CscartResponse) {
	for work := range jobs {
		fmt.Printf("Started worker #%d\n", id)
		body := prepareBody(work)
		//fmt.Printf("%v+", body)
		query := url.Values{}
		query.Set("ajax_custom", "1")
		result := api.Api("POST", "products", body, query)
		go func() {
			for key := range work.Options {
				if key == "" {
					return
				}
				body := prepareOptions(work, result, key)
				query := url.Values{}
				query.Set("ajax_custom", "1")
				result := api.Api("POST", "options", body, query)
				results <- *result
			}
		}()
		results <- *result
	}
	wg.Done()
}

func prepareOptions(w Work, res *api.CscartResponse, name string) string {
	//var j map[string]interface{} = make(map[string]interface{})
	data, ok := res.Data.(map[string]interface{}) //uhh!
	if !ok {
		log.Fatal("error converting data")
	}

	options := make(map[string]interface{})
	options["option_type"] = "S"
	options["option_name"] = name
	options["product_id"] = data["product_id"].(float64)
	variants := map[string]map[string]string{}
	for k, v := range w.Options[name] {
		variants[strconv.Itoa(k)] = map[string]string{"variant_name": v}
	}
	options["variants"] = variants
	body, err := json.Marshal(options)
	if err != nil {
		log.Fatal(err)
	}
	return string(body)
}

func prepareBody(work Work) string {
	var params map[string]interface{} = make(map[string]interface{})
	Sprice := strings.Trim(work.Amount, " ")
	price, _ := strconv.ParseFloat(Sprice, 64)
	price = price * exchangeRate                 // convert to ugx shs
	params["price"] = (price + (price * margin)) // add the margin as a percentage of a product price
	params["company_id"] = vendor
	params["product"] = work.Product
	params["full_description"] = work.Description
	params["list_price"] = price * 1
	params["amount"] = 20
	params["main_category"] = category
	params["category_ids"] = []string{category}
	params["product_features"] = map[string]map[string]string{"285": {"feature_type": "T", "value": "Huddah"}}
	if len(work.Imagepairs) > 0 {
		params["main_pair"] = map[string]map[string]string{"detailed": {"http_image_path": work.Imagepairs[0], "image_path": work.Imagepairs[0]}}
	}
	var pairs []map[string]map[string]string = make([]map[string]map[string]string, 0)
	for i := 1; i <= len(work.Imagepairs)-1; i++ {
		pair := map[string]map[string]string{"detailed": {"http_image_path": work.Imagepairs[i], "image_path": work.Imagepairs[i]}}
		pairs = append(pairs, pair)
	}
	params["image_pairs"] = pairs

	jsonbody, err := json.Marshal(params)
	if err != nil {
		log.Fatal(err)
	}
	return string(jsonbody)
}
