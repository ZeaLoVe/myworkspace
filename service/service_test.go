package service

import (
	"log"
	"testing"
)

func TestServiceDump(t *testing.T) {
	var s Service
	s.SetDefault()
	passtest := true
	if s.Name != "defaultservice" {
		passtest = false
	}
	if s.Port != 8080 {
		passtest = false
	}
	if s.Weight != 100 {
		passtest = false
	}
	if s.Priority != 20 {
		passtest = false
	}
	if s.Ttl != 10 {
		passtest = false
	}
	if s.Text != "default text for record something" {
		passtest = false
	}
	if passtest != true {
		t.Fatalf("setDefault error")
	}
}

func TestParseJSON(t *testing.T) {
	var ser Service
	ser.SetDefault()
	if res, err := ser.ParseJSON(); err == nil {
		t.Log(string(res))
		log.Println("test service parseJSON success")
	} else {
		t.Fatalf("test parseJSON fail")
	}
}

func TestLoadConfigFile(t *testing.T) {
	var ser Service
	ser.LoadConfigFile("config.json")
	ser.Dump()
}

func TestSetMachines(t *testing.T) {
	var ser Service
	ser.setMachines(nil)
	ser.Dump()
}
