package service

import (
	"errors"
	"github.com/coreos/go-etcd/etcd"
	"log"
	"path"
	"strings"
)

type Register struct {
	s *Service
}

func (r *Register) ChangeMachines(newMachine []string) {
	if len(newMachine) == 0 {
		log.Println("No etcd Machine address given")
		return
	} else {
		r.s.machines = newMachine
		log.Println("etcd Machine address changed success")
	}
}

func (r *Register) UpdateService() error {

	r.s.setHost("")
	tmpList := strings.Split(r.s.Key, ".")
	for i, j := 0, len(tmpList)-1; i < j; i, j = i+1, j-1 {
		tmpList[i], tmpList[j] = tmpList[j], tmpList[i]
	}

	key := path.Join(append([]string{"/skydns/"}, tmpList...)...)
	value, err := r.s.ParseJSON()
	if err != nil {
		log.Printf("can't get value in UpdateService")
		return err
	}
	log.Printf("#UPdateService#insert key: %v\n", key)
	log.Printf("#UPdateService#insert value: %v\n", string(value))

	if len(r.s.machines) == 0 {
		log.Fatalf("No etcd machines")
		return errors.New("No etcd machines")
	}
	client := etcd.NewClient(r.s.machines)

	_, errSet := client.Set(key, string(value), r.s.Ttl)
	if errSet != nil {
		return err
	} else {
		return nil
	}
}
