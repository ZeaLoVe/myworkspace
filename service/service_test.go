package service

import (
	"fmt"
	"testing"
)

func TestServiceDump(t *testing.T) {
	var ser Service
	ser.SetDefault()
	ser.Dump()
}

func TestParseJSON(t *testing.T) {
	var ser Service
	ser.SetDefault()
	if res, err := ser.ParseJSON(); err == nil {
		fmt.Println(string(res))
	} else {
		t.Fatalf("test parseJSON fail")
	}
}

func TestLoadConfigFile(t *testing.T) {
	var ser Service
	ser.SetDefault()
	ser.LoadConfigFile("config.json")
	ser.Dump()
}
