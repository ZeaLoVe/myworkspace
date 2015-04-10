package service

import (
	"fmt"
	"testing"
)

func TestHealthCheckDump(t *testing.T) {
	var hc HealthCheck
	hc.SetDefault()
	hc.Dump()
}

func TestHealthCheck(t *testing.T) {
	var hc HealthCheck
	hc.SetDefault()
	if res, err := hc.Check(); err != nil {
		t.Fatalf(err.Error())
	} else {
		if res != PASS {
			t.Fatalf("default health check fail")
		} else {
			fmt.Println("Script health check pass")
		}
	}

	hc.TTL = 0
	hc.Script = "something script"
	hc.Interval = 10
	if res, err := hc.Check(); err != nil {
		t.Fatalf(err.Error())
	} else {
		if res != PASS {
			t.Fatalf("Script health check fail")
		} else {
			fmt.Println("Script health check pass")
		}
	}

	hc.TTL = 0
	hc.HTTP = "http://baidu.com"
	hc.Interval = 10
	if res, err := hc.Check(); err != nil {
		t.Fatalf(err.Error())
	} else {
		if res != PASS {
			t.Fatalf("HTTP health check fail")
		} else {
			fmt.Println("HTTP health check pass")
		}
	}

}
