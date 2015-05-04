package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	. "sdagent/backends"
	. "sdagent/util"
	"strconv"
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

//priority is higher with smaller value
func (parser *ServiceParser) ReducePriority() {
	parser.Priority = parser.Priority + 10
}

func (parser *ServiceParser) IncreaseWeight(num uint64) {
	parser.Weight = parser.Weight + num
}

func (parser *ServiceParser) ReduceWeight(num uint64) {
	parser.Weight = parser.Weight - num
	if parser.Weight < 0 {
		parser.Weight = 100
	}
}

func (parser *ServiceParser) ToJSON() ([]byte, error) {
	return json.Marshal(parser)
}

//use to update DNS data in etcd
type Service struct {
	//etcd machines
	Node     string `json:"node,omitempty"`
	Name     string `json:"name,omitempty"`
	Host     string `json:"host,omitempty"`
	Port     uint64 `json:"port,omitempty"`
	Priority uint64 `json:"priority,omitempty"`
	Weight   uint64 `json:"weight,omitempty"`
	Text     string `json:"text,omitempty"`
	Ttl      uint64 `json:"ttl,omitempty"`
	Key      string `json:"key,omitempty"`

	Hc []HealthCheck `json:"checks,omitempty"`

	backend Backend `json:"-"`
}

func NewService() *Service {
	ser := new(Service)
	ser.SetDefault()
	return ser
}

func (s *Service) SetKey(key string) {

	if key == "" { //s.Key not set and given key is empty
		if s.Key != "" { //s.Key already set
			return
		}
		if s.Name == "" { //name is must
			log.Println("[WARN]Service Name not set")
			return
		}
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
	if host == "" {
		if s.Host != "" {
			return
		}
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
	if err := s.backend.SetMachines(newMachine); err != nil {
		log.Printf("[WARN]Set Machine Error, err:%v\n", err.Error())
	}
}

func (s *Service) DefaultServiceParser() *ServiceParser {
	var parser ServiceParser
	parser.Host = s.Host
	parser.Port = s.Port
	parser.Priority = s.Priority
	parser.Weight = s.Weight
	parser.Text = s.Text
	parser.Ttl = s.Ttl
	return &parser
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

//Service can't run job if Key,ttl not set
func (s *Service) CanRun() bool {
	if s.Key == "" || s.Ttl == 0 {
		return false
	} else {
		return true
	}
}

func (s *Service) CheckAll() int {
	if len(s.Hc) == 0 {
		log.Printf("[WARN]No health check in service: %v, ignore health check.\n", s.Key)
		return PASS
	}
	res := PASS
	for _, health := range s.Hc {
		oneres, err := health.Check()
		if err != nil && oneres == FAIL {
			return oneres
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

// A service need to call InitService before UpdateService,one time enough
func (s *Service) UpdateService(parser *ServiceParser) error {
	if s.Key == "" || s.Host == "" {
		return fmt.Errorf("Miss Key and Host")
	}
	key := GenKey(s.Key)

	var err error
	var value []byte
	if parser == nil {
		value, err = s.DefaultServiceParser().ToJSON()
	} else {
		value, err = parser.ToJSON()
	}

	if err != nil {
		return fmt.Errorf("Can't get value in UpdateService")
	}
	//log.Printf("[DEBUG]UPdateService key: %v.\n", key)
	//log.Printf("[DEBUG]UPdateService value: %v.\n", string(value))

	if err := s.backend.UpdateKV(key, string(value), s.Ttl); err == nil {
		return nil
	} else {
		return err
	}
}

func (s *Service) InitService() {
	s.SetKey("")
	s.SetHost("")
	s.SetMachines(nil)
}

//for init service
func (s *Service) SetDefault() {

	if s.Port <= 0 {
		s.Port = 8080
	}
	if s.Weight <= 0 {
		s.Weight = 100
	}
	if s.Priority <= 0 {
		s.Priority = 20
	}
	if s.Ttl <= 1 {
		s.Ttl = 10
	}
	if s.Text == "" {
		s.Text = "default text"
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
	fmt.Printf("%v health check set. \n", len(s.Hc))
	for _, health := range s.Hc {
		health.Dump()
	}
}
