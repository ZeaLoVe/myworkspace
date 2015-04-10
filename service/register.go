package service

import (
	"errors"
	"fmt"
	"github.com/coreos/go-etcd/etcd"
	"log"
	"path"
	"strings"
)

const ttloffset = 0 //

type Register struct {
	s *Service
}

func (r *Register) ChangeMachines(newMachine []string) {
	if len(newMachine) == 0 {
		log.Println("No Machine address")
		return
	} else {
		r.s.machines = newMachine
		log.Println("etcd Machine address changed success")
	}
}

func (r *Register) UpdateService() error {
	r.s.SetKey("")

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
	fmt.Printf("insert key: %v\n", key)
	fmt.Printf("insert value: %v\n", string(value))

	if len(r.s.machines) == 0 {
		log.Fatalf("No etcd machines")
		return errors.New("No etcd machines")
	}
	client := etcd.NewClient(r.s.machines)

	_, errSet := client.Set(key, string(value), (r.s.Ttl + ttloffset))
	if errSet != nil {
		return err
	} else {
		return nil
	}
}
