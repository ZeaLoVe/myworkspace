package sdagent

import (
	"fmt"
	"log"
	. "sdagent/service"
	"time"
)

const NORUNNINGNUM = 43
const KEEPALIVENUM = 34
const STOPCHANNUM = 10

const (
	PREPARE = iota
	READY
	RUNNING
)

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
}

func (state *JobState) SetSuccess() {
	state.UpdateCount++
}

func (state *JobState) SetWarn() {
	state.WarnCount++
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

func (j *Job) LastCheckState() int {
	return j.state.LastCheckStatus
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
		//update time must smaller than TTL, here make it smaller 1,to be considering...
		j.config.UpdateInterval = time.Duration(j.S.Ttl-1) * time.Second
	} else {
		return fmt.Errorf("No enough infomation for job setconfig")
	}
	j.stopChan = make(chan uint64)
	j.keepAliveChan = make(chan uint64, 100)
	j.SetJobState(READY) //PREPARE->READY , can run
	j.state.SetJobName(j.config.JobID)
	return nil
}

func (j *Job) jobStop() {
	j.SetJobState(PREPARE)
	close(j.stopChan)
	close(j.keepAliveChan)
}

func (j *Job) Run() {
	if j.config.JOBSTATE != READY {
		log.Printf("[WARM]jobID:%v state is not READY, run job fail.\n", j.config.JobID)
		return
	}
	j.SetJobState(RUNNING)
	defer j.jobStop()

	internal := time.After(0)
	timeout := time.After(j.config.UpdateInterval * 3)
	heartbeatSender := time.Tick(j.config.UpdateInterval)
	heartbeatReciever := time.Tick(j.config.UpdateInterval)
	timeoutcount := 0
	//check job's heartbeat to make sure it is alive

	go func() {
		for {
			select {
			case <-timeout:
				log.Printf("[ERR]jobID:%v timeout.\n", j.config.JobID)
				j.SetConfig()
				go j.Run()
				log.Printf("[INFO]jobID:%v restart.\n", j.config.JobID)
				return
			case <-heartbeatReciever:
				if keep, ok := <-j.keepAliveChan; ok {
					if keep == KEEPALIVENUM {
						timeout = time.After(j.config.UpdateInterval * 3) //reflesh timeout
					}
					if keep == NORUNNINGNUM {
						return
					}
				} else {
					return
				}
			}
			time.Sleep(1 * time.Second)
		}
	}()

	for {
		select {
		case <-internal: //do update by last check result
			internal = time.After(j.config.UpdateInterval)
			res := j.LastCheckState() //Get last check result
			if res == PASS {
				if timeoutcount >= 2 {
					log.Println("[WARN]timeout reset etcd machines")
					j.S.SetMachines(nil)
					timeoutcount = 0
				}
				if err := j.S.UpdateService(nil); err != nil {
					if err.Error() == "etcd timeout" {
						timeoutcount++
						continue
					}
					j.state.SetFail()
					log.Printf("[WARN]jobID:%v do updateservice fail,error:%v", j.config.JobID, err.Error())
				} else {
					j.state.SetSuccess()
					//log.Printf("[INFO]jobID:%v do updateservice success", j.config.JobID)
				}
			} else if res == WARN {
				j.state.SetWarn()
				log.Printf("[WARN]jobID:%v do health check warn", j.config.JobID)
			} else if res == FAIL {
				j.state.SetFail()
				log.Printf("[WARN]jobID:%v do health check Fail", j.config.JobID)
			} else {
				//nothing
			}
		case <-j.stopChan:
			j.keepAliveChan <- NORUNNINGNUM
			log.Printf("[INFO]jobID:%v stop\n", j.config.JobID)
			return
		case <-heartbeatSender:
			j.keepAliveChan <- KEEPALIVENUM
			j.state.IncHeartBeat()
			go func() {
				res := j.S.CheckAll()
				j.state.LastCheckStatus = res
			}()
		}
	}
}
