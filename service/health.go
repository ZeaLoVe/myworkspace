package service

import (
	"encoding/json"
	"fmt"
)

//Result of health Check
const (
	PASS = iota + 1
	WARN
	FAIL
)

type HealthCheck struct {
	CheckName string `json:"name,omitempty"`
	CheckID   string `json:"id,omitempty"`
	TTL       uint64 `json:"ttl,omitempty"`
	Script    string `json:"script,omitempty"`
	HTTP      string `json:"http,omitempty"`
	Interval  uint64 `json:"interval,omitempty"`
	Timeout   uint64 `json:"timeout,omitempty"`
	Notes     string `json:"notes,omitempty"`
}

func (hc *HealthCheck) SetDefault() {
	if hc.CheckName == "" {
		hc.CheckName = "defaultcheck"
	}
	if hc.CheckID == "" {
		hc.CheckID = "defaultid"
	}
	if hc.TTL == 0 {
		hc.TTL = 10
	}
	if hc.Interval == 0 {
		hc.Interval = 0
	}
	if hc.Notes == "" {
		hc.Notes = "this is default health check."
	}
}

func (hc *HealthCheck) ParseJSON() ([]byte, error) {
	return json.Marshal(hc)
}

func (hc *HealthCheck) Check() (int, error) {
	if hc.TTL != 0 { //TTL check
		return PASS, nil
	}
	if hc.Script != "" && hc.Interval != 0 {
		return PASS, nil
	}
	if hc.HTTP != "" && hc.Interval != 0 {
		return PASS, nil
	}
	return FAIL, nil
}

//for test ,dump value of health check.
func (hc *HealthCheck) Dump() {
	fmt.Printf("name:%v\n", hc.CheckName)
	fmt.Printf("id:%v\n", hc.CheckID)
	fmt.Printf("ttl:%v\n", hc.TTL)
	fmt.Printf("script:%v\n", hc.Script)
	fmt.Printf("http:%v\n", hc.HTTP)
	fmt.Printf("interval:%v\n", hc.Interval)
	fmt.Printf("timeout:%v\n", hc.Timeout)
	fmt.Printf("notes:%v\n", hc.Notes)
}
