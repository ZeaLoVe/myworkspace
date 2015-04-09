package service

import (
	"log"
)

type Service struct {
	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
	Priority int    `json:"priority,omitempty"`
	Weight   int    `json:"weight,omitempty"`
	Text     string `json:"text,omitempty"`
	Ttl      uint32 `json:"ttl,omitempty"`
	// etcd key where we found this service and ignore from json un-/marshalling
	Key string        `json:"-"`
	hc  []HealthCheck `json:"-"`
}
