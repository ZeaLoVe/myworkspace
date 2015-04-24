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
	JOBSTATE int //Runtime state of job: PREPARE->READY->RUNNING->PREPARE

	UpdateInterval time.Duration
	JobID          string //come from Service.key ,what is unique
}

//for statistic
type JobState struct {
	JobName         string `json:"JobName,omitempty"`
	FailCount       uint64 `json:"FailCount,omitempty"`
	UpdateCount     uint64 `json:"UpdateCount,omitempty"`
	WarnCount       uint64 `json:"WarnCount,omitempty"`
	HeartBeatSent   uint64 `json:"HeartBeatSent,omitempty"`
	LastCheckStatus int    `json:"LastCheckStatus,omitempty"`
}

func (state *JobState) SetJobName(name string) {
	state.JobName = name
}

func (state *JobState) SetFail() {
	state.FailCount++
	state.LastCheckStatus = FAIL
}

func (state *JobState) SetSuccess() {
	state.UpdateCount++
	state.LastCheckStatus = PASS
}

func (state *JobState) SetWarn() {
	state.WarnCount++
	state.LastCheckStatus = WARN
}

func (state *JobState) IncHeartBeat() {
	state.HeartBeatSent++
}

type Job struct {
	S             Service
	keepAliveChan chan uint64 //check whether job is down
	stopChan      chan uint64 //get cmd to stop
	config        JobConfig
	state         JobState
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
	j.state.SetJobName(j.config.JobID)
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
				if err := j.S.UpdateService(); err != nil {
					if err.Error() == "No etcd machines" {
						log.Printf("[ERR]jobID:%v No etcd machines.\n", j.config.JobID)
						continue
					}
					j.state.SetFail()
					log.Printf("[WARN]jobID:%v do updateservice fail,error:%v", j.config.JobID, err.Error())
				} else {
					j.state.SetSuccess()
					//log.Printf("[DEBUG]jobID:%v do updateservice success", j.config.JobID)
				}
			} else if res == WARN {
				j.state.SetWarn()
				log.Printf("[WARN]jobID:%v do health check Warn", j.config.JobID)
			} else if res == FAIL {
				j.state.SetFail()
				log.Printf("[WARN]jobID:%v do health check Fail", j.config.JobID)
			} else {
				//nothing
			}
		case <-j.stopChan:
			log.Printf("[DEBUG]jobID:%v stop.\n", j.config.JobID)
			return
		case <-heartbeat:
			j.keepAliveChan <- KEEPALIVENUM //no meanning
			j.state.IncHeartBeat()
		}
	}
}
