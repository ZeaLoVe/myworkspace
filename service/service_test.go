package service

import (
	"testing"
)

func TestLoadConfigFile(t *testing.T) {
	var ser Service
	//ser.SetDefault()
	ser.LoadConfigFile("config.json")
	ser.Dump()
}
