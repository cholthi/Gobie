package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

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

	closeHandler()
	var category string

	flag.StringVar(&category, "category", "", "Product category to scrap on Jumia")
	flag.Parse()

	categorycollector := colly.NewCollector(colly.UserAgent("Mozilla/5.0 (X11; Linux i586; rv:63.0) Gecko/20100101 Firefox/63.0"))
	productcollector := colly.NewCollector(colly.UserAgent("Mozilla/5.0 (X11; Linux i586; rv:63.0) Gecko/20100101 Firefox/63.0"))

	//lets scrap them category pages

	categorycollector.OnHTML("section.products div.sku", func(e *colly.HTMLElement) {
		link := e.ChildAttr("a", "href")
		title := e.ChildText("a h2.title")
		fmt.Printf("product Info: %s : %s", title, link)
		productcollector.Visit(link)
	})

	categorycollector.OnHTML("a[href] ", func(e *colly.HTMLElement) {

		if e.Attr("title") == "Next" {
			fmt.Printf("found Next Link: %s", e.Attr("href"))
			e.Request.Visit(e.Attr("href"))
		}
	})

	productcollector.OnHTML("main.osh-container", func(e *colly.HTMLElement) {
		product := product{}
		product.Name = e.ChildText("section.sku-detail div.details-wrapper div.details span h1.title")
		product.Price = e.ChildAttr("section.sku-detail div.details-wrapper div.details div.details-footer div.price-box div span.price span:nth-of-type(2)", "data-price")
		e.ForEach("section.sku-detail div.media div.thumbs-wrapper div#thumbs-slide a", func(index int, l *colly.HTMLElement) {
			src := l.Attr("href")
			product.ImagePairs = append(product.ImagePairs, src)
		})
		//Options for products
		option := make(map[string][]string)
		option_name := e.ChildText("section.sku-detail div.details-wrapper div.details div.detail-features div.sizes div.list div.title")
		option[option_name] = make([]string, 0)
		e.ForEach("section.sku-detail div.details-wrapper div.details div.detail-features div.sizes div.list span.sku-size", func(i int, e *colly.HTMLElement) {
			option[option_name] = append(option[option_name], e.Text)
		})
		product.Options = option
		desc := e.DOM.Find("section div.osh-tabs div#productDescriptionTab div.product-description")
		html, _ := desc.Html()
		product.Desc = strings.Replace(html, "Jumia", "Agoro", -1)
		products = append(products, product)
	})

	categorycollector.Visit("" + category)
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
