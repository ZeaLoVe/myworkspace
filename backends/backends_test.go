package backends

import (
	"log"
	"sdagent/util"
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

func TestCheckKV(t *testing.T) {
	var backend Backend
	if flag, err := backend.CheckKV(GenKey("172.24.133.22:8080")); !flag {
		log.Printf("Check return false with error: %v", err.Error())
	} else {
		log.Println("172.24.133.22:8080 check ok")
	}
	if flag, err := backend.CheckKV(GenKey("172.24.133.22:8103")); !flag {
		log.Printf("Check return false with error: %v", err.Error())
	} else {
		log.Println("172.24.133.22:8103 check ok")
	}
	if flag, err := backend.CheckKV(GenKey("192.168.181.16:8080")); !flag {
		log.Printf("Check return false with error: %v", err.Error())
	} else {
		log.Println("192.168.181.16:8080 check ok")
	}
}
