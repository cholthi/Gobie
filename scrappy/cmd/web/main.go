package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	scrapper "github.com/cholthi/scrappy/scrape"
)

type ScrapeRequest struct {
	Host      string `json:"host"`
	Category  string `json:"category"`
	NoRequest int    `json:"request_no"`
	Replace   string `json:"replace_text,omitempty"`
	FileName string  `json::file_name"`
}

type Response struct {
	Success bool `json:"success"`
	Number  int  `json:"number"`
}

var logger log.Logger

func initLogging() {
	file := "./agoro.log"
	output, err := os.OpenFile(file, os.O_WRONLY|os.O_APPEND, 0666)
	//defer output.Close()
	if err != nil {
		log.Println(err)
	}

	logger.SetOutput(output)
	logger.SetPrefix("topup-api")
	logger.SetFlags(log.Lshortfile | log.Ldate)
}

/**
Main attaches handlers to the router
and starts the server on static port
*/
func main() {
	initLogging()
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("."))))
	mux.HandleFunc("/home", http.HandlerFunc(home))
	mux.HandleFunc("/calculator", http.HandlerFunc(cal))
	mux.HandleFunc("/ajax/scrape", http.HandlerFunc(scrape))
	mux.HandleFunc("/ajax/upload", http.HandlerFunc(uploadToserver))

	log.Fatal(http.ListenAndServe(":8082", mux))
}

func ParseTemplates(rootDir string) (*template.Template, error) {
	cleanedRootDir := filepath.Clean(rootDir)
	pfx := len(cleanedRootDir) + 1
	rootTemplate := template.New("")

	err := filepath.Walk(cleanedRootDir, func(path string, info os.FileInfo, err1 error) error {
		if !info.IsDir() && strings.HasSuffix(path, ".tpl") {
			if err1 != nil {
				return err1
			}

			b, err2 := ioutil.ReadFile(path)
			if err2 != nil {
				return err2
			}
			templateName := path[pfx:]
			t := rootTemplate.New(templateName)
			_, err2 = t.Parse(string(b))
			if err2 != nil {
				return err2
			}
		}
		return nil
	})
	return rootTemplate, err
}

func home(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tpl := template.Must(template.New("home.pl").ParseFiles("templates/home.tpl"))
		/*if err != nil {
			http.Error(rw, http.StatusText(http.StatusInternalServerError), 500)
			return
		}*/
		err := tpl.ExecuteTemplate(rw, "home.tpl", nil)
		if err != nil {
			panic(err)
		}
		return
	}
}

func cal(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tpl := template.Must(template.New("cal.pl").ParseFiles("templates/cal.tpl"))
		/*if err != nil {
			http.Error(rw, http.StatusText(http.StatusInternalServerError), 500)
			return
		}*/
		err := tpl.ExecuteTemplate(rw, "cal.tpl", nil)
		if err != nil {
			panic(err)
		}
		return
	}
}

func scrape(rw http.ResponseWriter, r *http.Request) {
	hosts := map[string]string{"Uganda": "https://www.jumia.ug", "Kenya": "https://www.jumia.co.ke"}
	req := ScrapeRequest{}
	r.ParseForm()
	country := r.FormValue("host")
	host := hosts[country]
	req.Host = host
	req.Category = r.FormValue("category")
	req.FileName = r.FormValue("file")
	req.Replace = r.FormValue("replace")
	no, err := strconv.Atoi(r.FormValue("request_no"))
	if err != nil {
		no = 100
	}
	req.NoRequest = no

	//args := []string{"--category", req.Category, "--host", req.Host, "--file", req.File, "--replace", req.Replace, "--request-no", r.FormValue("request_no")}
	item, err := reqToItem(req)
	if err != nil {
		logger.Print(err)
		res := []byte(`{"success": false, "number":0}`)
		rw.Header().Set("Content-Type", "application/json")
		rw.Write(res)
		return
	}

	products := scrapper.Scrape(*item)
	logger.Print(len(products))
	if len(products) >= 2 {
		err := products.Encode(item.Inventory)
		if err != nil {
			logger.Print(err)
			res := []byte(`{"success": false, "number":0}`)
			rw.Header().Set("Content-Type", "application/json")
			rw.Write(res)
			return
		}
		l := len(products)
		//sl := strconv.Itoa(l)
		data := fmt.Sprintf(`{"success": true, "number":%d, "file":%q}`, l, item.Inventory.Name())
		logger.Println(item.Inventory.Name())
		item.Inventory.Close()
		rw.Header().Set("Content-Type", "application/json")
		rw.Write([]byte(data))
		return
	}

	res := []byte(`{"success": false, "number":0}`)
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(res)
	return
}

func uploadToserver(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	category := r.FormValue("category")
	vendor := r.FormValue("vendor")
	margin := r.FormValue("margin")
	rate := r.FormValue("rate")
	file := r.FormValue("file")
	args := []string{"--category", category, "--vendor", vendor, "--file", file, "--rate", rate, "--margin", margin}
	cmd := exec.Command("./cscart", args...)
	cmd.Stdout = logger.Writer()

	err := cmd.Run()
	if err != nil {
		logger.Print(err)
		res := []byte(`{"success": false, "number":0}`)
		rw.Header().Set("Content-Type", "application/json")
		rw.Write(res)
		return
	}
	res := []byte(`{"success": true, "number":0}`)
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(res)
	return
}

func reqToItem(req ScrapeRequest) (*scrapper.ScrapeItem, error) {
	ret := &scrapper.ScrapeItem{}
	u, err := url.Parse(req.Host)
	if err != nil {
		return nil, err
	}
	ret.Host = u
	ret.Category = req.Category

	fobj, err := os.OpenFile(req.FileName, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return nil, err
	}
	ret.Inventory = fobj
	ret.NoRequest = req.NoRequest
	ret.Replace = req.Replace

	return ret, nil
}
