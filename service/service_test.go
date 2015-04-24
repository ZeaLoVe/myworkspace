package service

import (
	"log"
	"testing"
)

func (s *Service) SetTestMachines() {
	s.machines = []string{"http://192.168.181.16:2379"}
}

func TestServiceDump(t *testing.T) {
	var s Service
	s.SetDefault()
	passtest := true
	if s.Name != "default" {
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
	s.InitService()
	if s.Key != (s.Node + ".default") {
		t.Fatalf("set Key error in initservice")
	}
	if s.Host != "192.168.48.110" {
		t.Fatalf("set Host error in initservice")
	}
	//if s.machines[0] != "http://192.168.181.16:2379" {
	//	t.Fatalf("set machines error in initservice")
	//}
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

func TestSetMachines(t *testing.T) {
	var ser Service
	ser.SetMachines(nil)
	//if ser.machines[0] != "http://192.168.181.16:2379" {
	//	t.Fatalf("SetMachine nill fail")
	//}
	tmp := []string{"http://127.0.0.1:2379"}
	ser.SetMachines(tmp)
	if ser.machines[0] != "http://127.0.0.1:2379" {
		t.Fatalf("SetMachine given address fail")
	}
	ser.SetMachines(nil)
	if ser.machines[0] != "http://127.0.0.1:2379" {
		t.Fatalf("SetMachine already set given empty fail")
	}
}

func TestSetHost(t *testing.T) {
	service := NewService()
	service.SetHost("")
	if service.Host != "192.168.48.110" {
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
}

func TestSetKey(t *testing.T) {
	service := NewService()
	service.Node = "s1"
	service.SetKey("")
	if service.Key != "s1.default" {
		t.Fatal("Set Key default error")
	}
	service.SetKey("hsp.sd.dp")
	log.Println(service.Key)
	if service.Key != "hsp.sd.dp" {
		t.Fatal("Set Key by given name error")
	}
	service.SetKey("")
	if service.Key != "hsp.sd.dp" {
		t.Fatal("Set Key empty error")
	}
}

func TestCheckAll(t *testing.T) {
	ser := NewService()
	ser.LoadConfigFile("config.json")
	ser.SetDefault()
	ser.InitService()
	if ser.CheckAll() == PASS {
		t.Fatal("Check default can't pass")
	}
	ser.Hc[0].Script = "ipconfig"
	if ser.CheckAll() != PASS {
		t.Fatal("Check ALL not pass")
	}
	ser.Hc[0].Script = "exit 0"
	if ser.CheckAll() != PASS {
		t.Fatal("Check ALL not pass")
	}
	ser.Hc[0].Script = "exit 2"
	if ser.CheckAll() != FAIL {
		t.Fatal("Check ALL not pass")
	}
	ser.Hc[0].Script = "exit 1"
	if ser.CheckAll() != WARN {
		t.Fatal("Check ALL not pass")
	}
	ser.Hc[1].HTTP = "www.google.nd"
	if ser.CheckAll() != FAIL {
		t.Fatal("Check ALL not pass")
	}
	ser.Hc = nil
	if ser.CheckAll() != PASS {
		t.Fatal("Check ALL Empty check not pass")
	}
}

func TestUpdateService(t *testing.T) {
	ser := NewService()
	ser.LoadConfigFile("config.json")
	ser.SetDefault()
	ser.InitService()
	err := ser.UpdateService()
	if err != nil {
		t.Fatal("update service error")
	} else {
		log.Printf("update service pass")
	}
	//t.Fatal("see log for")
}
