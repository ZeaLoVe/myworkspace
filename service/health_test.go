package service

import (
	"log"
	"testing"
)

func TestHealthCheckDump(t *testing.T) {
	var hc HealthCheck
	hc.SetDefault()
	passtest := true
	if hc.CheckName == "" {
		passtest = false
	}
	if hc.CheckID == "" {
		passtest = false
	}
	if hc.Timeout != 10 {
		passtest = false
	}
	if hc.Interval != 10 {
		passtest = false
	}
	if hc.Notes != "Health check Notes not given." {
		passtest = false
	}
	if passtest != true {
		t.Fatalf("SetDefault error")
	} else {
		log.Println("Healthcheck Setdefault success")
	}
}

func TestHealthCheckParseJSON(t *testing.T) {
	hc := NewHealthCheck()
	if res, err := hc.ParseJSON(); err == nil {
		//t.Log(string(res))
		log.Println("test healthCheck parseJSON success")
	} else {
		t.Fatalf("test healthCheck parseJSON fail got '%v'\n", res)
	}
}

func TestHealthCheck(t *testing.T) {
	hc := NewHealthCheck()
	if res, err := hc.Check(); err != nil {
		log.Println(err.Error())
	} else {
		if res != PASS {
			t.Fatalf("default health check fail")
		} else {
			log.Println("Script health check pass")
		}
	}

	hc.TTL = 0
	hc.Script = ""
	hc.HTTP = "http://baidu.com"
	if res, err := hc.Check(); err != nil {
		t.Fatalf(err.Error())
	} else {
		if res != PASS {
			t.Fatalf("HTTP health check fail")
		} else {
			log.Println("HTTP health check pass")
		}
	}

	hc.HTTP = "https://google.nd"
	if res, err := hc.Check(); err == nil {
		t.Fatalf("no error http check here")
	} else {
		if res == PASS {
			t.Fatalf("Fail http check return pass")
		} else {
			log.Println("HTTP health check fail check rights")
		}
	}

	hc.TTL = 0
	hc.Script = "ipconfig/all"
	if res, err := hc.Check(); err != nil {
		t.Fatalf(err.Error())
	} else {
		if res != PASS {
			t.Fatalf("Script health check fail")
		} else {
			log.Println("Script health check pass")
		}
	}

	hc.Script = "dd -t"
	if res, err := hc.Check(); err == nil {
		t.Fatalf("no error script check here")
	} else {
		if res == PASS {
			t.Fatalf("Fail script check return pass")
		} else {
			log.Println("Script health check fail check right")
		}
	}

}
