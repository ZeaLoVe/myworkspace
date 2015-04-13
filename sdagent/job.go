package sdagent

import (
	"fmt"
	"log"
	. "myworkspace/service"
	"time"
)

type JobConfig struct {
	LastCheckStatus int
	UpdateInterval  time.Duration
	JobID           string //come from Service.key ,what is unique
}

type Job struct {
	r             *Register
	s             *Service
	keepAliveChan chan uint64 //check whether job is down
	stopChan      chan uint64 //get cmd to stop
	config        JobConfig
}

func (j *Job) SetConfig() {
	if j.s != nil {
		j.config.JobID = j.s.Key
		j.config.UpdateInterval = time.Duration(j.s.Ttl/2) * time.Second //update time must smaller than TTL
	}
}

func (j *Job) Run() {
	j.config.LastCheckStatus = PASS
	internal := time.After(j.config.UpdateInterval)
	log.Printf("Job run interval:%v\n", j.config.UpdateInterval)
	for {
		select {
		case <-internal:
			if j.config.LastCheckStatus == PASS {
				fmt.Println("job update")
				j.r.UpdateService()
			}
		case <-j.stopChan:
			fmt.Println("job stop")
		}
	}
}
