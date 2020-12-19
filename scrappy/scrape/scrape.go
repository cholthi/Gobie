package scrape

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

type product struct {
	Name       string              `json:"product"`
	Price      string              `json:"price"`
	Desc       string              `json:"full_description"`
	ImagePairs []string            `json:"image_pairs"`
	Options    map[string][]string `json:"options"`
}

type ScrapeItem struct {
	Host      *url.URL
	NoRequest int
	Category  string
	Replace   string   //replace the instant of this string in the description
	Inventory *os.File // file to write the scrapped data to
}

type Products []product

func Scrape(item ScrapeItem) Products {
	count := 0
	products := Products{}
	ua := getUserAgent()
	categorycollector := colly.NewCollector(colly.UserAgent(ua))
	categorycollector.Limit(&colly.LimitRule{
		DomainGlob:  "*jumia.*",
		Parallelism: 3,
		//Delay:      5 * time.Second,
	})
	//categorycollector.SetProxy("https://103.215.157.53:58338/")
	productcollector := colly.NewCollector(colly.UserAgent(ua))
	//productcollector.SetProxy("https://103.192.38.103:8082/")
	productcollector.Limit(&colly.LimitRule{
		DomainGlob:  "*jumia.*",
		Parallelism: 3,
		//Delay:      5 * time.Second,
	})

	//lets scrap 'em category pages

	categorycollector.OnHTML("div.row article.c-prd", func(e *colly.HTMLElement) {
		link := e.ChildAttr("a.core", "href")
		title := e.ChildText("a div.info h3.name")
		fmt.Printf("product Info: %s : %s", title, item.Host.Host+link)
		productcollector.Visit("https://" + item.Host.Host + link) //make full path
	})

	categorycollector.OnHTML("a[href] ", func(e *colly.HTMLElement) {

		if e.Attr("aria-label") == "Next Page" {
			fmt.Printf("found Next Link: %s", e.Attr("href"))
			e.Request.Visit(e.Attr("href"))
		}
	})

	productcollector.OnHTML("main.-pvs div.row section.col12", func(e *colly.HTMLElement) {

		if count > item.NoRequest {
			return
		}
		product := product{}
		product.Name = e.ChildText("div.row div.col10 div.-df div.-fs0 h1.-pts")
		replacers := strings.Split(item.Replace, "|")
		price := e.ChildText("div.row div.col10 div.-phs div.-mtxs span.-fs24")
		if contains := strings.Contains(price, replacers[0]); contains {
			price = strings.Replace(price, replacers[0], "", -1)
			price = strings.Replace(price, ",", "", -1)
			price = strings.Trim(price, " ")
			if isrange := strings.Contains(price, "-"); isrange {
				price = strings.Split(price, "-")[0]
				price = strings.Trim(price, " ")
			}
		}
		product.Price = price
		e.ForEach("div.row div.col6 div.-pbs div#imgs a", func(index int, l *colly.HTMLElement) {
			src := l.Attr("href")
			product.ImagePairs = append(product.ImagePairs, src)
		})
		//Options for products
		option := make(map[string][]string, 0)
		option_name := e.ChildText("section.col12 div.row div.col10 div.-phxs div.-mhxs span")
		option[option_name] = make([]string, 0)
		e.ForEach("section.col12 div.row div.col10 div.-phxs div.var-w label", func(i int, e *colly.HTMLElement) {
			option[option_name] = append(option[option_name], e.Text)
		})

		product.Options = option
		parent := e.DOM.ParentsUntil("div.row main.-pvs")
		html, _ := parent.Find("div.row div.col12 div.card div.markup").Html()
		//html, _ := desc.Html()
		product.Desc = strings.ReplaceAll(html, "Jumia", "Agoro")
		product.Desc = strings.ReplaceAll(html, replacers[1], "South Sudan")
		products = append(products, product)
		count++
	})

	split := strings.Split(item.Category, "?")
	if split[0] != "" {
		item.Host.Path = split[0]
	}
	if len(split) == 2 {
		item.Host.RawQuery = split[1]
	}

	fmt.Println(item.Host.String())
	categorycollector.Visit(item.Host.String()) //starts scraping
	return products
}

func (p Products) Encode(file io.WriteCloser) error {

	enc := json.NewEncoder(file)
	enc.SetEscapeHTML(false)
	enc.SetIndent("	", "")
	err := enc.Encode(p) //here
	if err != nil {
		return err
	}
	return nil
}
