package sdagent

import (
	"errors"
	"log"
	. "myworkspace/service"
	"time"
)

const KEEPALIVENUM = 34
const STOPCHANNUM = 10

const (
	PREPARE = iota
	READY
	RUNNING
)

//const HeartBeatInterval = time.Duration(5 * time.Second) //heartbeat time

type JobConfig struct {
	JOBSTATE        int //Runtime state of job: PREPARE->READY->RUNNING->PREPARE
	LastCheckStatus int //keep for multithread healthcheck
	UpdateInterval  time.Duration
	JobID           string //come from Service.key ,what is unique
}

type Job struct {
	S             Service
	keepAliveChan chan uint64 //check whether job is down
	stopChan      chan uint64 //get cmd to stop
	config        JobConfig
}

func NewJob() *Job {
	job := new(Job)
	return job
}

func (j *Job) SetJobState(state int) {
	j.config.JOBSTATE = state
}

func (j *Job) CanRun() bool {
	return j.S.CanRun()
}

func (j *Job) SetConfig() error {
	if j.S.Key != "" && j.S.Ttl != 0 {
		j.config.JobID = j.S.Key
		//update time must smaller than TTL, to be considering...
		j.config.UpdateInterval = time.Duration(j.S.Ttl-1) * time.Second
	} else {
		log.Println("[WARM]No enough infomation for job SetConfig")
		return errors.New("No enough infomation for job setconfig")
	}
	j.stopChan = make(chan uint64)
	j.keepAliveChan = make(chan uint64, 10)
	j.SetJobState(READY) //PREPARE->READY , can run
	return nil
}

func (j *Job) jobStop() {
	close(j.stopChan)
	close(j.keepAliveChan)
	j.SetJobState(PREPARE)
}

func (j *Job) Run() {
	if j.config.JOBSTATE != READY {
		log.Printf("[WARM]jobID:%v state is not READY, run job fail.\n", j.config.JobID)
		return
	}
	j.SetJobState(RUNNING)
	defer j.jobStop()
	internal := time.Tick(j.config.UpdateInterval)
	heartbeat := time.Tick(j.config.UpdateInterval / 2)
	log.Printf("[DEBUG]JobID: %v run interval:%v\n", j.config.JobID, j.config.UpdateInterval)
	for {
		select {
		case <-internal: //do update and check
			res := j.S.CheckAll()
			if res == PASS {
				log.Printf("[DEBUG]jobID:%v update service Success", j.config.JobID)
				if err := j.S.UpdateService(); err != nil {
					if err.Error() == "No etcd machines" {
						log.Printf("[ERR]jobID:%v No etcd machines.\n", j.config.JobID)
						time.Sleep(time.Second * 1)
					}
					log.Printf("[WARM]jobID:%v do updateservice fail,error:%v", j.config.JobID, err.Error())
				} else {
					log.Printf("[DEBUG]jobID:%v do updateservice success", j.config.JobID)
				}
			} else if res == WARN {
				log.Printf("[WARM]jobID:%v do health check Warn", j.config.JobID)
			} else if res == FAIL {
				log.Printf("[WARM]jobID:%v do health check Fail", j.config.JobID)
			} else {
				//nothing
			}
		case <-j.stopChan:
			log.Printf("[DEBUG]jobID:%v stop.\n", j.config.JobID)
			return
		case <-heartbeat:
			j.keepAliveChan <- KEEPALIVENUM //no meanning
			//log.Printf("[DEBUG]jobID:%v heartbeat recieve!", j.config.JobID)
		}
	}
}
