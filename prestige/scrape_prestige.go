package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gocolly/colly"
)

type product struct {
	Name       string              `json:"product"`
	Price      string              `json:"price"`
	Desc       string              `json:"full_description"`
	ImagePairs []string            `json:"image_pairs"`
	Options    map[string][]string `json:"options"`
}

var products []product = []product{}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	closeHandler()
	var category string

	flag.StringVar(&category, "category", "", "Product category to scrap on Prestige")
	flag.Parse()

	ua := getUserAgent()
	categorycollector := colly.NewCollector(colly.UserAgent(ua))
	//categorycollector.SetProxy("https://103.215.157.53:58338/")
	productcollector := colly.NewCollector(colly.UserAgent(ua))
	//productcollector.SetProxy("https://103.192.38.103:8082/")

	//lets scrap them category pages

	categorycollector.OnHTML("div.products section", func(e *colly.HTMLElement) {
		link := e.ChildAttr("div.product-inner div.product-thumbnail-wrapper a:last-child", "href")
		//title := e.ChildText("section.product div.product-meta-wrapper div.prod-meta-relative div.prod-meta-relative-wrapper div.prod-meta-relative-inner h3.product-title a")
		title := e.ChildText("div.product-inner div.product-meta-wrapper div.prod-meta-relative div.prod-meta-relative-wrapper div.prod-meta-relative-inner h3.product-title a")
		fmt.Printf("product Info:Link %s : Title %s", link, title)
		productcollector.Visit(link)
	})

	categorycollector.OnHTML("li a.next ", func(e *colly.HTMLElement) {

		//if e.Attr("class") == "next" {
		fmt.Printf("found Next Link: %s", e.Attr("href"))
		e.Request.Visit(e.Attr("href"))
		//}
	})

	productcollector.OnHTML("div#content div.product", func(e *colly.HTMLElement) {
		product := product{}
		product.Name = e.ChildText("div.product-image-summary-wrapper div.summary h1")
		price := e.ChildText("div.product-image-summary-wrapper div.summary div.woocommerce-product-box-wrapper p.price span")
		if ugx := strings.Contains(price, "KES"); ugx {
			price = strings.Replace(price, "KES", "", -1)
			price = strings.Trim(price, "\xc2\xa0")
			price = strings.Replace(price, "Kenya", "South Sudan", -1)
			price = strings.Replace(price, "prestige", "Agoro", -1)
			price = strings.Replace(price, ",", "", -1)
			//price = strings.Replace(price, "&nbsp;", "", -1)
			if isrange := strings.Contains(price, "-"); isrange {
				price = strings.Split(price, "-")[0]
				price = strings.Trim(price, " ")
			}
		}
		product.Price = price
		e.ForEach("div.product-image-summary-wrapper div.images div.p_image", func(index int, l *colly.HTMLElement) {
			src := l.ChildAttr("div.item a", "href")
			product.ImagePairs = append(product.ImagePairs, src)
		})
		//Options for products
		option := make(map[string][]string, 0)
		option_name := "Variation"
		option[option_name] = make([]string, 0)
		option[option_name] = append(option[option_name], "Paper Bag")

		product.Options = option
		//parent := e.DOM.ParentsUntil("div.row main.-pvs")
		html, _ := e.DOM.Find("div.product-image-summary-wrapper div.summary div:nth-of-type(1)").Html()
		//html, _ := desc.Html()
		product.Desc = strings.Replace(html, "Prestige", "Agoro", -1)
		product.Desc = strings.Replace(html, "Kenya", "South Sudan", -1)
		products = append(products, product)
	})

	categorycollector.Visit("https://prestigebookshop.com" + category)
	file, err := os.OpenFile("scrapped.json", os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	enc := json.NewEncoder(file)
	enc.SetEscapeHTML(false)
	enc.SetIndent("	", "")
	enc.Encode(products)
	fmt.Printf("Scrapped %d products", len(products)+1)
}

func closeHandler() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		file, err := os.OpenFile("scrapped.json", os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			panic(err)
		}
		enc := json.NewEncoder(file)
		enc.SetEscapeHTML(false)
		enc.SetIndent("	", "")
		enc.Encode(products)
		fmt.Println()
		fmt.Printf("Scrapped %d products", len(products))
		fmt.Println()
		fmt.Println("Shuting Down now...")
		os.Exit(0)
	}()
}

func getUserAgent() string {
	ua := []string{
		"Mozilla/5.0 (X11; U; Linux Core i7-4980HQ; de; rv:32.0; compatible; JobboerseBot; http://www.jobboerse.com/bot.htm) Gecko/20100101 Firefox/38.0",
		"Apache/2.4.34 (Ubuntu) OpenSSL/1.1.1 (internal dummy connection)",
		"Mozilla/5.0 (X11; U; Linux Core i7-4980HQ; de; rv:32.0; compatible; JobboerseBot; https://www.jobboerse.com/bot.htm) Gecko/20100101 Firefox/38.0",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) HeadlessChrome/74.0.3729.157 Safari/537.36",
		"Mozilla/5.0 (X11; U; Linux i686; pt-BR; rv:1.9.0.15) Gecko/2009102815 Ubuntu/9.04 (jaunty) Firefox/3.0.15",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/44.0.2403.157 Safari/537.36",
	}

	randomidx := rand.Intn(len(ua) - 1)
	return ua[randomidx]
}
