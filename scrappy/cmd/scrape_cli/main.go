package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"github.com/cholthi/scrappy/scrape"
)

func main() {
	var category string
	var replacer string
	var host string
	var inventory string
	var noRequests int

	flag.StringVar(&category, "category", "", "Product category to scrap on Jumia")
	flag.StringVar(&replacer, "replace", "UGX|Uganda", "string to replace in the scrapped body")
	flag.StringVar(&host, "host", "https://jumia.ug", "the host scrape")
	flag.StringVar(&inventory, "file", "./scrapped.json", "the file to store the scrapped data")
	flag.IntVar(&noRequests, "request-no", 0, "No of scrapped requests pages to process")
	flag.Parse()

	item := scrape.ScrapeItem{}
	item.Category = category
	item.NoRequest = noRequests

	file, err := os.OpenFile(inventory, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	item.Inventory = file
	item.Replace = replacer

	u, err := url.Parse(host)
	if err != nil {
		panic(err)
	}

	item.Host = u

	prducts := scrape.Scrape(item)
	closeHandler(item, prducts)
	err = prducts.Encode(item.Inventory)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Scrapped %d products", len(prducts)+1)
}

func closeHandler(item scrape.ScrapeItem, p scrape.Products) {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		/*file, err := os.OpenFile(item.Inventory, os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			panic(err)
		}*/
		p.Encode(item.Inventory)
		fmt.Printf("Scrapped %d products", len(p))
		fmt.Println()
		fmt.Println("Shuting Down now...")
		os.Exit(0)
	}()
}
