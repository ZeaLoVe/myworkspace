package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"

	"myworkspace/util"
)

const ETCDPORT = "2379"
const ETCDMACHINES = "etcd.sdp"

//get ip by name ,use for etcd machines discoury
func getipByName(name string) []string {
	ns, err := net.LookupIP(name)
	if err != nil {
		fmt.Printf("no ips for %v", name)
		return nil
	} else {
		var ips []string
		for _, ip := range ns {
			ips = append(ips, ip.String())
		}
		return ips
	}
}

type ServiceParser struct {
	Host     string `json:"host,omitempty"`
	Port     uint64 `json:"port,omitempty"`
	Priority uint64 `json:"priority,omitempty"`
	Weight   uint64 `json:"weight,omitempty"`
	Text     string `json:"text,omitempty"`
	Ttl      uint64 `json:"ttl,omitempty"`
}

//use to update DNS data in etcd
type Service struct {
	//etcd machines

	Machines string `json:"machines,omitempty"`
	Node     string `json:"node,omitempty"`
	Name     string `json:"name,omitempty"`
	Host     string `json:"host,omitempty"`
	Port     uint64 `json:"port,omitempty"`
	Priority uint64 `json:"priority,omitempty"`
	Weight   uint64 `json:"weight,omitempty"`
	Text     string `json:"text,omitempty"`
	Ttl      uint64 `json:"ttl,omitempty"`
	Key      string `json:"key,omitempty"`
	//use to do healthcheck
	Hc []HealthCheck `json:"checks"`

	machines []string `json:"-"`
}

func (s *Service) SetKey(key string) {
	if s.Key != "" { //s.Key already set
		return
	}
	if key == "" { //s.Key not set and given key is empty
		if s.Node != "" {
			s.Key = s.Node + "." + s.Name
		} else { //if node not set ,given a randan node name
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			s.Node = "s" + strconv.Itoa(r.Intn(10000))
			s.Key = s.Node + "." + s.Name
		}
	} else { //s.Key not set and given key is not empty
		s.Key = key
	}
}

func (s *Service) setHost(host string) {
	if s.Host != "" {
		return
	}
	if host == "" {
		pip, err := util.GetPrivateIP()
		if err != nil {
			log.Fatal("Get PrivateIP error")
		} else {
			s.Host = pip.String()
		}

	} else {
		s.Host = host
	}
}

func (s *Service) setMachines(newMachine []string) {
	if len(newMachine) == 0 {
		if len(s.machines) != 0 {
			return //already set
		}
		if len(s.machines) == 0 && s.Machines != "" {
			s.machines = strings.Split(s.Machines, ",")
		} else {
			var tmpMachines []string
			tmpMachines = getipByName(ETCDMACHINES)
			for i, machine := range tmpMachines {
				tmpMachines[i] = "http://" + machine + ":" + ETCDPORT
			}
			s.machines = tmpMachines
		}
	} else {
		s.machines = newMachine
	}
}

func (s *Service) ParseJSON() ([]byte, error) {
	var parser ServiceParser
	parser.Host = s.Host
	parser.Port = s.Port
	parser.Priority = s.Priority
	parser.Weight = s.Weight
	parser.Text = s.Text
	parser.Ttl = s.Ttl
	return json.Marshal(parser)
}

func (s *Service) LoadConfigFile(filename string) {
	config, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal("Can't load config file,can't get infomation of service")
	}
	if err := json.Unmarshal(config, s); err != nil {
		log.Fatal("Unmarsh to service struct fail,can't get infomation of service")
	}
	s.SetDefault()
}

func (s *Service) CheckAll() int {
	res := PASS
	for _, health := range s.Hc {
		oneres, err := health.Check()
		if err != nil {
			return FAIL
		} else {
			if oneres == WARN || oneres == FAIL {
				res = oneres
				return res
			}
		}
	}
	return res
}

//for test
func (s *Service) SetDefault() {
	if s.Name == "" {
		s.Name = "defaultservice"
	}
	if s.Port == 0 {
		s.Port = 8080
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
	s.setMachines(nil)
}

//for test
func (s *Service) Dump() {
	fmt.Printf("key:%v\n", s.Key)
	fmt.Printf("name:%v\n", s.Name)
	fmt.Printf("node:%v\n", s.Node)
	fmt.Printf("text:%v\n", s.Text)
	fmt.Printf("host:%v\n", s.Host)
	fmt.Printf("port:%v\n", s.Port)
	fmt.Printf("priority:%v\n", s.Priority)
	fmt.Printf("weight:%v\n", s.Weight)
	fmt.Printf("ttl:%v\n", s.Ttl)
	fmt.Printf("machines:%v\n", s.machines)
	fmt.Printf("%v health check set. \n", len(s.Hc))
	for _, health := range s.Hc {
		health.Dump()
	}
}
