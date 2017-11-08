package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello world")
}
func handlerEvaluationTrigger(w http.ResponseWriter, r *http.Request) {
	forceInvokeClient()
}

func handlerNewHook(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var ticket Ticket

	err := decoder.Decode(&ticket)
	if err != nil {
		panic(err)
	}

	defer r.Body.Close()

	ticket.ID = bson.NewObjectId()
	insertData("tickets", ticket)
	fmt.Fprintln(w, ticket.ID.Hex())
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 Ok"))
}

func handlerAccessHook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")

	var ticket Ticket

	objectid := bson.ObjectIdHex(vars["id"])

	c := session.DB("CurrencyDB").C("tickets")
	c.FindId(objectid).One(&ticket)

	// if ticket.id == "" {
	//
	// }

	resp, err := json.MarshalIndent(ticket, "", "   ")
	if err != nil {
		panic(err)
	}

	fmt.Fprint(w, string(resp))
}

func handlerDeleteHook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")

	objectid := bson.ObjectIdHex(vars["id"])

	c := session.DB("CurrencyDB").C("tickets")
	c.RemoveId(objectid)
}

func handlerLatest(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var ticket Ticket
	var rate CurrencyData

	err := decoder.Decode(&ticket)
	if err != nil {
		panic(err)
	}

	defer r.Body.Close()

	c := session.DB("CurrencyDB").C("rates")

	c.Find(bson.M{}).Sort("-_id").One(&rate)
	fmt.Fprintln(w, rate.Rates[ticket.Target])
}

func handlerAverage(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var ticket Ticket
	var rates []CurrencyData
	var avg float64
	err := decoder.Decode(&ticket)
	if err != nil {
		panic(err)
	}

	defer r.Body.Close()

	c := session.DB("CurrencyDB").C("rates")
	c.Find(bson.M{}).Sort("-_id").Limit(7).All(&rates)

	for _, val := range rates {
		avg += val.Rates[ticket.Target]
	}

	avg /= float64(len(rates))

	fmt.Fprintln(w, avg)
}
