package main

import "gopkg.in/mgo.v2/bson"

// CurrencyData struct
type CurrencyData struct {
	ID    bson.ObjectId      `bson:"_id,omitempty"`
	Base  string             `json:"base"`
	Date  string             `json:"date"`
	Rates map[string]float64 `json:"rates"`
}

// Ticket struct
type Ticket struct {
	ID         bson.ObjectId `bson:"_id,omitempty"`
	URL        string        `json:"webhookURL"`
	Base       string        `json:"baseCurrency"`
	Target     string        `json:"targetCurrency"`
	Rate       float64       `json:"currentRate"`
	MinTrigger float64       `json:"minTriggerValue"`
	MaxTrigger float64       `json:"maxTriggerValue"`
}

// InvokedData struct
type InvokedData struct {
	Base       string  `json:"baseCurrency"`
	Target     string  `json:"targetCurrency"`
	Rate       float64 `json:"currentRate"`
	MinTrigger float64 `json:"minTriggerValue"`
	MaxTrigger float64 `json:"maxTriggerValue"`
}
