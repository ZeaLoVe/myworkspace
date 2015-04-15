package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os/exec"
	"strconv"
	"syscall"
	"time"

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
	HTTP      string `json:"http,omitempty"`     //just for health check type ,value has no meanning
	Interval  uint64 `json:"interval,omitempty"` //no use now,keep it
	Timeout   uint64 `json:"timeout,omitempty"`  //no use now,keep it
	Notes     string `json:"notes,omitempty"`    //no use now,keep it
}

func NewHealthCheck() *HealthCheck {
	hc := new(HealthCheck)
	hc.SetDefault()
	return hc
}

func (hc *HealthCheck) SetDefault() {
	if hc.CheckName == "" {
		hc.CheckName = "chk_name"
	}
	if hc.CheckID == "" {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		hc.CheckID = "chk_id" + strconv.Itoa(r.Intn(100000))
	}
	if hc.Timeout == 0 {
		hc.Timeout = 10
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

// timeout must be set
func (hc *HealthCheck) ScriptCheck() (int, error) {
	cmd, err := util.ExecScript(hc.Script)
	if err != nil {
		log.Printf("[WARM]Fail to setup invoke '%v' with err:'%v'.\n", hc.Script, err.Error())
		return FAIL, err
	}
	//output, err := cmd.Output()
	//log.Printf("Script return: %v ", string(output))

	if err := cmd.Start(); err != nil {
		log.Printf("[WARM]Fail to invoke '%v' with err:'%v'.\n", hc.Script, err.Error())
		return FAIL, err
	}
	err = cmd.Wait() // get cmd return value
	if err == nil {
		return PASS, nil
	}
	exitErr, ok := err.(*exec.ExitError)
	if ok {
		if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
			code := status.ExitStatus()
			log.Printf("[DEBUG]Script Check:'%v' return: %v .\n", hc.Script, code)
			if code == 0 {
				log.Printf("[DEBUG]Script Check:'%v' is passing.\n", hc.Script)
				return PASS, nil
			} else if code == 1 {
				log.Printf("[DEBUG]Script Check:'%v' is warning.\n", hc.Script)
				return WARN, err
			} else {
				log.Printf("[DEBUG]Script Check:'%v' is failing.\n", hc.Script)
				return FAIL, err
			}
		}
	}
	return FAIL, err
}

func (hc *HealthCheck) HttpCheck() (int, error) {
	resp, err := http.Get(hc.HTTP)
	if err != nil {
		log.Printf("[WARN]Http request:'%v' failed. error:'%v'.\n", hc.HTTP, err)
		return FAIL, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[WARN]Get '%v' error while reading http body:'%v'.\n", err, body)
	}
	//log.Printf("http request:'%v' get status: '%v' with body:'%s'\n", hc.HTTP, resp.Status, body)
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		return PASS, nil
	} else if resp.StatusCode == 429 {
		return WARN, errors.New("Too many querys")
	} else {
		return FAIL, errors.New("Http check not pass")
	}
	return PASS, nil
}

func (hc *HealthCheck) Check() (int, error) {
	//return success while not set,but warm will be logged
	if hc.TTL == 0 && hc.Script == "" && hc.HTTP == "" {
		log.Printf("[WARM]Health check config miss.\n")
		return PASS, errors.New("miss health check config")
	}
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
