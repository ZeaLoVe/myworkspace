package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
)

//use to update DNS data in etcd
type Service struct {
	//etcd machines
	machines []string `json:"etcd,machine,omitempty"`
	Node     string   `json:"node,omitempty"`
	Name     string   `json:"name,omitempty"`
	Host     string   `json:"host,omitempty"`
	Port     uint64   `json:"port,omitempty"`
	Priority uint64   `json:"priority,omitempty"`
	Weight   uint64   `json:"weight,omitempty"`
	Text     string   `json:"text,omitempty"`
	Ttl      uint64   `json:"ttl,omitempty"`
	// etcd key where we found this service and ignore from json (un)marshalling
	Key string `json:"key,omitempty"`

	//use to do healthcheck
	hc []HealthCheck `json:"checks"`
}

func (s *Service) SetKey(key string) {
	if s.Key != "" { //s.Key already set
		return
	}
	if key == "" { //s.Key not set and given key is empty
		s.Key = strconv.Itoa(int(s.Port)) + "." + s.Node + "." + s.Name
		fmt.Printf("SetKEY:-------------------------: %v\n", s.Key)
	} else { //s.Key not set and given key is not empty
		s.Key = key
	}
}

//for test
func (s *Service) setHost(host string) {
	if s.Host != "" {
		return
	}
	if host == "" {
		s.Host = "192.168.1.1" //here need be changed by function getCurrent IP
	} else {
		s.Host = host
	}
}

func (s *Service) ParseJSON() ([]byte, error) {
	return json.Marshal(s)
}

func (s *Service) LoadConfigFile(filename string) {
	config, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal("Can't load config file,can't get infomation of service")
	}
	if err := json.Unmarshal(config, s); err != nil {
		log.Fatal("Unmarsh to service struct fail,can't get infomation of service")
	}
}

//for test
func (s *Service) SetDefault() {
	if len(s.machines) == 0 {
		s.machines = []string{"http://192.168.181.16:2379"}
	}
	if s.Name == "" {
		s.Name = "defaultservice"
	}
	if s.Node == "" {
		s.Node = "localhost"
	}
	if s.Port == 0 {
		s.Port = 80
	}
	if s.Weight == 0 {
		s.Weight = 100
	}
	if s.Priority == 0 {
		s.Priority = 20
	}
	if s.Ttl == 0 {
		s.Ttl = 10
	}
	if s.Text == "" {
		s.Text = "default text for record something"
	}
	s.SetKey("")
	s.setHost("")
}

//for test
func (s *Service) Dump() {
	for _, str := range s.machines {
		fmt.Printf("etcd Machines:%v\n", str)
	}
	fmt.Printf("key:%v\n", s.Key)
	fmt.Printf("name:%v\n", s.Name)
	fmt.Printf("host:%v\n", s.Host)
	fmt.Printf("port:%v\n", s.Port)
	fmt.Printf("priority:%v\n", s.Priority)
	fmt.Printf("weight:%v\n", s.Weight)
	fmt.Printf("ttl:%v\n", s.Ttl)
	fmt.Printf("machines:%v\n", s.machines)
	if len(s.hc) != 0 {
		for _, health := range s.hc {
			health.Dump()
		}
	}
}
