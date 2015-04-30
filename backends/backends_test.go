package backends

import (
	"log"
	"myworkspace/util"
	"testing"
)

func TestGenKey(t *testing.T) {
	log.Println(GenKey("my.com"))
	log.Println(GenKey("sdp.cn"))
	log.Println(GenKey("hdm.com.dn.cn"))
	log.Println(GenKey("16.1212.121"))
	log.Println(GenKey("998.11.11"))
}

func TestSetMachines(t *testing.T) {
	var backend Backend
	tmpList := util.GetIPByName("zealove.xyz")
	backend.SetMachines(tmpList)
	if backend.client == nil {
		log.Println("setmachines fail")
	} else {
		log.Println("Setmaachines succes")
	}
}

func TestUpdateKV(t *testing.T) {
	var backend Backend
	if err := backend.UpdateService("/skydns/cn/nd", "try it", 10); err == nil {
		log.Println("Update success")
	} else {
		t.Fatalf("Update fail error:%v", err.Error())
	}
	tmpList := util.GetIPByName("baidu.com")
	for _, machine := range tmpList {
		log.Println(machine)
	}
	backend.SetMachines(tmpList)
	if err := backend.UpdateService("/skydns/nddd/ddd", "{something:value}", 10); err == nil {
		t.Fatal("Update success")
	} else {
		log.Printf("Update fail error:%v", err.Error())
	}
}
