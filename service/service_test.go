package service

import (
	"log"
	"testing"
)

func TestServiceDump(t *testing.T) {
	var s Service
	s.SetDefault()
	passtest := true
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
	if passtest != true {
		t.Fatalf("setDefault error")
	}
	s.InitService()

	if s.Host != "192.168.252.44" {
		t.Fatalf("set Host error in initservice")
	}
	//if s.machines[0] != "http://192.168.181.16:2379" {
	//	t.Fatalf("set machines error in initservice")
	//}
}

func TestParseJSON(t *testing.T) {
	var ser Service
	ser.SetDefault()
	if res, err := ser.DefaultServiceParser().ToJSON(); err == nil {
		t.Log(string(res))
		log.Println("test service parseJSON success")
	} else {
		t.Fatalf("test parseJSON fail")
	}
}

func TestSetHost(t *testing.T) {
	service := NewService()
	service.SetHost("")
	if service.Host != "192.168.252.44" {
		t.Fatalf("Set Host by private ip error")
	} //本机IP
	service.SetHost("local.host")
	if service.Host != "local.host" {
		t.Fatal("Set Host by given name error")
	}
	service.SetHost("")
	if service.Host != "local.host" {
		t.Fatal("Set Host by given empty error")
	}
	service.SetHost("eXample.sdp.xxxX")
	if service.Host != "example.sdp.xxxx" {
		t.Fatal("set Host tolower fail")
	}
}

func TestSetKey(t *testing.T) {
	service := NewService()
	service.Name = "eXampLe.sdp"
	service.Node = "S1"
	service.SetKey("")
	if service.Key != "s1.example.sdp" {
		t.Fatal("Set Key default error")
	}
	service.SetKey("Hsp.sd.dP")
	log.Println(service.Key)
	if service.Key != "hsp.sd.dp" {
		t.Fatal("Set Key by given name error")
	}
	service.SetKey("")
	if service.Key != "hsp.sd.dp" {
		t.Fatal("Set Key empty error")
	}
}

//func TestCheckAll(t *testing.T) {
//	ser := NewService()
//	ser.LoadConfigFile("config.json")
//	ser.SetDefault()
//	ser.InitService()
//	if ser.CheckAll() == PASS {
//		t.Fatal("Check default can't pass")
//	}
//	ser.Hc[0].Script = "ipconfig"
//	if ser.CheckAll() != PASS {
//		t.Fatal("Check ALL not pass")
//	}
//	ser.Hc[0].Script = "exit 0"
//	if ser.CheckAll() != PASS {
//		t.Fatal("Check ALL not pass")
//	}
//	ser.Hc[0].Script = "exit 2"
//	if ser.CheckAll() != FAIL {
//		t.Fatal("Check ALL not pass")
//	}
//	ser.Hc[0].Script = "exit 1"
//	if ser.CheckAll() != WARN {
//		t.Fatal("Check ALL not pass")
//	}
//	ser.Hc[1].HTTP = "www.google.nd"
//	if ser.CheckAll() != FAIL {
//		t.Fatal("Check ALL not pass")
//	}
//	ser.Hc = nil
//	if ser.CheckAll() != PASS {
//		t.Fatal("Check ALL Empty check not pass")
//	}
//}

func TestCanRun(t *testing.T) {
	ser := NewService()
	ser.LoadConfigFile("config.json")
	if ser.CanRun() {
		t.Fatal("empty key can run")
	}
	ser.Ttl = 0
	if ser.CanRun() {
		t.Fatal("empty ttl can run")
	}
	ser.SetDefault()
	ser.InitService()
	if !ser.CanRun() {
		t.Fatal("can't run")
	}
	ser.SetHost("eXample.sdP")
	if !ser.CanRun() {
		t.Fatal("uperCase can't run")
	}
	ser.SetHost("falcon-ops.sdp.nd")
	if !ser.CanRun() {
		t.Fatal("common domain can't run")
	}
	ser.SetHost("ds_xxx.dsd.x")
	if ser.CanRun() {
		t.Fatal("invalid _ can run")
	}
	ser.SetHost("..ds.xdsd")
	if ser.CanRun() {
		t.Fatal("invalid .. can run")
	}
	ser.SetHost("dsdsx")
	if ser.CanRun() {
		t.Fatal("invalid domain length can run")
	}
}

func TestUpdateService(t *testing.T) {
	ser := NewService()
	ser.LoadConfigFile("config.json")
	//	ser.SetDefault()
	//	ser.InitService()
	//	err := ser.UpdateService(nil)
	//	if err != nil {
	//		t.Fatal("update service error")
	//	} else {
	//		log.Printf("update service pass")
	//	}

	if ser.isValidService() {
		log.Println("valid service")
	} else {
		log.Println("invalide service")
	}
	//t.Fatal("see log for")
}
