package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/coreos/go-etcd/etcd"
	"io/ioutil"
	"log"
	"math/rand"
	. "myworkspace/util"
	"path"
	"strconv"
	"strings"
	"time"
)

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

func NewService() *Service {
	ser := new(Service)
	ser.SetDefault()
	return ser
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
		pip, err := GetPrivateIP()
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
			tmpMachines := GetIPByName(ETCDMACHINES)
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

//it can't run if Key,ttl,machines not set
func (s *Service) CanRun() bool {
	if s.Key == "" || s.Ttl == 0 || len(s.machines) == 0 {
		return false
	} else {
		return true
	}
}

func (s *Service) CheckAll() int {
	if len(s.Hc) == 0 {
		log.Printf("No health check in service: %v, ignore health check.\n", s.Name)
		return PASS
	}
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

func (s *Service) UpdateService() error {

	s.setHost("")
	tmpList := strings.Split(s.Key, ".")
	for i, j := 0, len(tmpList)-1; i < j; i, j = i+1, j-1 {
		tmpList[i], tmpList[j] = tmpList[j], tmpList[i]
	}

	key := path.Join(append([]string{"/skydns/"}, tmpList...)...)
	value, err := s.ParseJSON()
	if err != nil {
		log.Printf("can't get value in UpdateService")
		return err
	}
	log.Printf("#UPdateService#insert key: %v\n", key)
	log.Printf("#UPdateService#insert value: %v\n", string(value))

	if len(s.machines) == 0 {
		log.Fatalf("No etcd machines")
		return errors.New("No etcd machines")
	}
	client := etcd.NewClient(s.machines)

	_, errSet := client.Set(key, string(value), s.Ttl)
	if errSet != nil {
		return err
	} else {
		return nil
	}
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
