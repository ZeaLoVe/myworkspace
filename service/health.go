package service

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"syscall"

	"myworkspace/util"
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
		hc.Interval = 10
	}
	if hc.Notes == "" {
		hc.Notes = "Health check Notes not given."
	}
}

func (hc *HealthCheck) ParseJSON() ([]byte, error) {
	return json.Marshal(hc)
}

func (hc *HealthCheck) TTLCheck() (int, error) {
	return PASS, nil
}

func (hc *HealthCheck) ScriptCheck() (int, error) {
	cmd, err := util.ExecScript(hc.Script)
	if err != nil {
		log.Printf("fail to setup invoke '%v' with err:'%v'", hc.Script, err.Error())
		return FAIL, err
	}
	output, err := cmd.Output()
	log.Printf("Script return:%v", string(output))
	exitErr, ok := err.(*exec.ExitError)
	if ok {
		if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
			code := status.ExitStatus()
			if code == 1 {
				log.Printf("check '%v' is warning ", hc.Script)
				return WARN, err
			}
		}
	}
	return PASS, err
}

func (hc *HealthCheck) HttpCheck() (int, error) {
	return PASS, nil
}

func (hc *HealthCheck) Check() (int, error) {
	if hc.TTL != 0 { //TTL check
		return hc.TTLCheck()
	}
	if hc.Script != "" && hc.Interval != 0 {
		return hc.ScriptCheck()
	}
	if hc.HTTP != "" && hc.Interval != 0 {
		return hc.HttpCheck()
	}
	return PASS, nil
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
