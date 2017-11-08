package main

import (
	"fmt"
	"net/http"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/mux"
	"github.com/robfig/cron"
)

var (
	session *mgo.Session
)

func setHandlers(r *mux.Router) {
	// Start http server
	r.HandleFunc("/", handlerNewHook).Methods("POST")
	r.HandleFunc("/latest", handlerLatest).Methods("POST")
	r.HandleFunc("/average", handlerAverage).Methods("POST")
	r.HandleFunc("/evaluationtrigger", handlerEvaluationTrigger).Methods("GET")
	r.HandleFunc("/{id}", handlerAccessHook).Methods("GET")
	r.HandleFunc("/{id}", handlerDeleteHook).Methods("DELETE")
}

func main() {
	startDataBase("127.0.0.1")
	getRates()

	// New mux router
	r := mux.NewRouter()
	setHandlers(r)

	fmt.Println("Server running on port 8080")
	http.ListenAndServe(":80", r)

	// Start cron job
	cron := cron.New()
	cron.Start()

	// Updates rates daily
	cron.AddFunc("@daily", updateAndInvoke)

	// Stop cron job
	cron.Stop()
	defer session.Close() // Close session
}

func updateAndInvoke() {
	updateRates()
	invokeClient()
}

func invokeClient() {
	var tickets []Ticket
	c := session.DB("CurrencyDB").C("tickets")
	c.Find(bson.M{}).All(&tickets)

	var rates CurrencyData
	c = session.DB("CurrencyDB").C("rates")
	c.Find(bson.M{}).Sort("-_id").One(&rates)

	for _, elem := range tickets {
		if outOfBounds(elem, rates) {
			notifyClient(elem, rates)
		}
	}
}

func forceInvokeClient() {
	var tickets []Ticket
	c := session.DB("CurrencyDB").C("tickets")
	c.Find(bson.M{}).All(&tickets)

	var rates CurrencyData
	c = session.DB("CurrencyDB").C("rates")
	c.Find(bson.M{}).Sort("-_id").One(&rates)

	for _, elem := range tickets {
		notifyClient(elem, rates)
	}
}
