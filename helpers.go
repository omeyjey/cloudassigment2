package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	mgo "gopkg.in/mgo.v2"
)

func startDataBase(url string) {
	var err error
	session, err = mgo.Dial(url)
	if err != nil {
		log.Fatal(err)
	}
}

func fetchRates(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatal(err)
	}
	resp.Body.Close()

	return body
}

func insertData(collection string, data interface{}) {
	c := session.DB("CurrencyDB").C(collection)

	// Insert Datas
	err := c.Insert(data)
	if err != nil {
		return
	}
}

func getRates() {
	body := fetchRates("http://api.fixer.io/latest?base=EUR")
	data := CurrencyData{}

	jsonError := json.Unmarshal(body, &data)
	if jsonError != nil {
		log.Fatal(jsonError)
	}

	insertData("rates", data)
}

func updateRates() {
	body := fetchRates("http://api.fixer.io/latest?base=EUR")

	data := CurrencyData{}

	jsonError := json.Unmarshal(body, &data)
	if jsonError != nil {
		log.Fatal(jsonError)
	}

	insertData("rates", data)
}

func outOfBounds(t Ticket, r CurrencyData) bool {
	return r.Rates[t.Target] > t.MaxTrigger || r.Rates[t.Target] < t.MinTrigger
}

func notifyClient(t Ticket, r CurrencyData) {
	invoke := InvokedData{
		Base:       t.Base,
		Target:     t.Target,
		Rate:       r.Rates[t.Target],
		MinTrigger: t.MinTrigger,
		MaxTrigger: t.MaxTrigger,
	}

	body, err := json.MarshalIndent(invoke, "", "   ")
	if err != nil {
		log.Fatal(err)
	}

	http.Post(t.URL, "application/x-wwww-form-urlencoded", bytes.NewBuffer(body))
}
