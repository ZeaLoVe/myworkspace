package service

import (
	"github.com/coreos/go-etcd/etcd"
	"log"
)

type Register struct {
	c        *etcd.client //etcd client
	machines []string
}

func (r *Register) Init() {

}

func (r *Register) DeleteService(s *Service) {

}

func (r *Register) UpdateService(s *Service) {

}
