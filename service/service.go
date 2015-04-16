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
	Host     string `json:"host,omitempty"` //need for DNS records
	Port     uint64 `json:"port,omitempty"`
	Priority uint64 `json:"priority,omitempty"`
	Weight   uint64 `json:"weight,omitempty"`
	Text     string `json:"text,omitempty"`
	Ttl      uint64 `json:"ttl,omitempty"` //need for DNS records ttl
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

	Hc []HealthCheck `json:"checks"`

	machines []string `json:"-"`
	client   *etcd.Client
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
			s.Node = "s" + strconv.Itoa(r.Intn(100000))
			s.Key = s.Node + "." + s.Name
		}
	} else { //s.Key not set and given key is not empty
		s.Key = key
	}
}

func (s *Service) SetHost(host string) {
	if s.Host != "" {
		return
	}
	if host == "" {
		ip, err := GetPrivateIP()
		if err != nil {
			log.Printf("[EER]Service of host not set and can't get private IP.\n")
		} else {
			s.Host = ip.String()
		}

	} else {
		s.Host = host
	}
}

func (s *Service) SetMachines(newMachine []string) {
	if len(newMachine) == 0 {
		if len(s.machines) != 0 {
			return //already set
		}
		if len(s.machines) == 0 && s.Machines != "" { //split etcd machines by ,
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

//this func no use in agent,just for test
func (s *Service) LoadConfigFile(filename string) {
	config, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal("Can't load config file,can't get infomation of service")
	}
	if err := json.Unmarshal(config, s); err != nil {
		log.Fatal("Unmarsh to service struct fail,can't get infomation of service")
	}
}

//Service can't run job if Key,ttl,machines not set
func (s *Service) CanRun() bool {
	if s.Key == "" || s.Ttl == 0 || len(s.machines) == 0 {
		return false
	} else {
		return true
	}
}

func (s *Service) CheckAll() int {
	if len(s.Hc) == 0 {
		log.Printf("[ERR]No health check in service: %v, ignore health check.\n", s.Key)
		return PASS
	}
	res := PASS
	for _, health := range s.Hc {
		oneres, err := health.Check()
		if err != nil {
			return FAIL
		} else {
			if oneres == FAIL {
				return FAIL
			}
			if oneres == WARN {
				res = WARN
			}
		}
	}
	return res
}

// call InitService before call this
func (s *Service) UpdateService() error {
	tmpList := strings.Split(s.Key, ".")
	for i, j := 0, len(tmpList)-1; i < j; i, j = i+1, j-1 {
		tmpList[i], tmpList[j] = tmpList[j], tmpList[i]
	}

	key := path.Join(append([]string{"/skydns/"}, tmpList...)...)
	value, err := s.ParseJSON()
	if err != nil {
		log.Printf("[WARM]Can't get value in function UpdateService.\n")
		return err
	}
	log.Printf("[DEBUG]UPdateService key: %v.\n", key)
	log.Printf("[DEBUG]UPdateService value: %v.\n", string(value))

	if len(s.machines) == 0 {
		log.Printf("[ERR]Service:%v No etcd machines.\n", s.Key)
		return errors.New("No etcd machines")
	}
	if s.client == nil {
		s.client = etcd.NewClient(s.machines)
	}

	// update first,then set
	_, errSet := s.client.Update(key, string(value), s.Ttl)
	if errSet == nil {
		return nil
	}
	_, errSet = s.client.Set(key, string(value), s.Ttl)
	if errSet != nil {
		return errSet
	} else {
		return nil
	}
}

func (s *Service) InitService() {
	s.SetKey("")
	s.SetHost("")
	s.SetMachines(nil)
}

//for init service
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
