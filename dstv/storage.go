package main

import (
	"encoding/csv"
	"os"
	"time"
)

//Store is a type that implements storage of subscriptions to packages
//implementations can allow to store in file system,remote or database system

type Store interface {
	Persist(RechargeRequest) error
}

type CsvStore struct {
	file *os.File
}

func (c CsvStore) Persist(pack RechargeRequest) error {
	csvw := csv.NewWriter(c.file)
	data := encode(pack) // returns []string which csv package encodes as csv row
	err := csvw.Write(data)
	if err != nil {
		logger.Println(err)
		return err
	}
	csvw.Flush()
	return nil
}

func NewCsvStore(f string) CsvStore {
	fo, err := os.OpenFile(f, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logger.Panic(err)
	}
	//defer fo.Close()
	store := CsvStore{file: fo}
	return store
}

func encode(pack RechargeRequest) []string {
	var out []string = make([]string, 0)

	//price := strconv.FormatFloat(float64(pack.), 'f', 6, 32)
	date := time.Now().Format(time.RFC3339)
	//format date,smartcardno,amount,transactionid,phoneNumber
	out = []string{date, pack.CardNumber, pack.Amount, pack.TxnId, pack.PhoneNo}
	return out
}
