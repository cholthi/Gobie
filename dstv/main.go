package main

import (
	"log"
	"net/http"
	"os"
	"time"
)

const (
	SMTP_HOST     = "smtp.gmail.com"
	SMTP_PORT     = 465
	SMTP_USERNAME = "swangin@dstvsouthsudan.com"
	SMTP_PASSWORD = "dstvsouthsudan2019"
)

var logfile string = "/root/dstv.log"
var csvdatabase string = "/root/MGURUSH_CUSTOMER_SUBSCRIBTION_UAT.csv"
var logger log.Logger = NewLogger(logfile)

var store Store = NewCsvStore(csvdatabase)

func NewLogger(file string) log.Logger {
	output, err := os.OpenFile(file, os.O_WRONLY|os.O_APPEND, 0666)
	//defer output.Close()
	if err != nil {
		log.Println(err)
	}

	logger := log.Logger{}
	logger.SetOutput(output)
	logger.SetPrefix("dstv")
	logger.SetFlags(log.Lshortfile)
	return logger
}

func main() {
	mux := http.DefaultServeMux
	mux.HandleFunc("/apiv1/packages", apiPackages)
	mux.HandleFunc("/apiv1/package", getPackage)
	mux.HandleFunc("/apiv1/recharge", apiRecharge)

	server := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		Handler:      authenticateMiddleware(mux),
		Addr:         ":8080",
		//TLSConfig:    configTLS(),
	}

	logger.Fatal(server.ListenAndServe())
}

/*func configTLS() *tls.Config {
	var config *tls.Config = new(tls.Config)
	CA_bundle, err := ioutil.ReadFile("/etc/ssl/certs/agoro_co.ca-bundle")
	if err != nil {
		logger.Panic(err)
	}
	CApool := x509.NewCertPool()
	CApool.AppendCertsFromPEM(CA_bundle)
	logger.Println("CA cert loaded")
	config.ClientCAs = CApool
	config.GetCertificate = func(info *tls.ClientHelloInfo) (certificate *tls.Certificate, e error) {
		c, err := tls.LoadX509KeyPair("/etc/ssl/certs/agoro_co.crt", "/etc/ssl/private/agoro.co.key")
		if err != nil {
			logger.Printf("Error loading key pair: %v\n", err)
			return nil, err
		}
		return &c, nil
	}
	return config
}*/
