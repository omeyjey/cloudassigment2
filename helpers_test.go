package main

import (
	"testing"
)

// Code from https://stackoverflow.com/questions/31595791/how-to-test-panics
// that shows how to test if a function panics

func TestStartDataBase(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("startDatabase panicked")
		}
	}()
	startDataBase("127.0.0.1")
}

func TestStartDataBaseShouldFail(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("StartDataBase did not panic")
		}
	}()
	startDataBase("mongodb://<dbuser>:<dbpassword>@ds041154.mongolab.com:41154/location")
}

func TestInsertDataShouldFail(t *testing.T) {
	var data interface{}
	data = 42

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("StartDataBase did not panic")
		}
	}()
	insertData("test", data)
}

func TestFetchRates(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("StartDataBase did panic")
		}
	}()
	fetchRates("http://api.fixer.io/latest?base=EUR")
}

func TestNotifyClient(t *testing.T) {
	ticket := Ticket{}
	rate := CurrencyData{}

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("NotifyClient did panic")
		}
	}()
	notifyClient(ticket, rate)
}
