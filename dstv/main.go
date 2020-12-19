package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
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
var csvdatabase string = "/root/MGURUSH_CUSTOMER_SUBSCRIBTION.csv"
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
	logger.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
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
		Addr:         ":4443",
		TLSConfig:    configTLS(),
	}

	logger.Fatal(server.ListenAndServeTLS("/etc/ssl/ssl.crt/agoro_co.crt", "/etc/ssl/ssl.key/agoro.key"))
}

func configTLS() *tls.Config {
	var config *tls.Config = new(tls.Config)
	CA_bundle, err := ioutil.ReadFile("/etc/ssl/ssl.crt/agoro_co.ca-bundle")
	if err != nil {
		logger.Panic(err)
	}
	CApool := x509.NewCertPool()
	if ok := CApool.AppendCertsFromPEM(CA_bundle); !ok {
		logger.Panicln("Failed to parse certificate authority file")
	}
	logger.Println("CA cert loaded")
	config.RootCAs = CApool
	return config
}
