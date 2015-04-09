package service

import (
	"fmt"
	"log"
)

//Result of health Check
const (
	PASS = iota + 1
	WARN
	FAIL
)

type HealthCheck struct {
	CheckName string
	CheckID   string
	TTL       uint64
	Script    string
	HTTP      string
	Timeout   uint64
	Notes     string
}

func (hc *HealthCheck) Check() (int, error) {
	if TTL == 0 { //TTL 监测
		return PASS, nil
	}
	if Script == "" {

	}
	if HTTP == "" {

	}
}
