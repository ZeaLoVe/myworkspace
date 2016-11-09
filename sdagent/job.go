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

		//update time must smaller than TTL, here make it half,to be considering...
		//UpdateInterval depends on ttl
		//if ttl > 5 min(300) than it's 300 , if ttl < 10 second(10) than it's 5, else make it half of ttl
		if j.S.Ttl >= 300 {
			j.config.UpdateInterval = 300 * time.Second
		} else if j.S.Ttl <= 10 {
			j.config.UpdateInterval = 5 * time.Second
		} else {
			j.config.UpdateInterval = time.Duration(j.S.Ttl/2) * time.Second
		}
	} else {
		return fmt.Errorf("No enough infomation for job setconfig")
	}
	j.stopChan = make(chan uint64)
	j.keepAliveChan = make(chan uint64, 1024)
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

	//init health check
	go func() {
		res := j.S.CheckAll()
		if res == WARN {
			log.Printf("[WARN]jobID:%v do health check warn", j.config.JobID)
		} else if res == FAIL {
			log.Printf("[WARN]jobID:%v do health check Fail", j.config.JobID)
		}
		j.state.LastCheckStatus = res
	}()

	interval := time.After(0)
	checkInterval := j.config.UpdateInterval
	timeout := time.After(checkInterval * 3)
	heartbeatSender := time.Tick(checkInterval)
	heartbeatReciever := time.Tick(checkInterval)
	//check job's heartbeat to make sure it is alive

	go func() {
		for {
			select {
			case <-timeout:
				log.Printf("[ERR]jobID:%v timeout.\n", j.config.JobID)
				return
			case <-heartbeatReciever:
				if keep, ok := <-j.keepAliveChan; ok {
					if keep == KEEPALIVENUM {
						timeout = time.After(checkInterval * 3) //reflesh timeout
					}
					if keep == NORUNNINGNUM {
						return
					}
				} else {
					return
				}
			}
			time.Sleep(10 * time.Millisecond)
		}
	}()

	for {
		select {
		case <-interval: //do update by last check result
			interval = time.After(j.config.UpdateInterval)
			res := j.LastCheckState() //Get last check result
			if res == PASS {
				go func() {
					//retry 2 times
					waitSecond := j.config.UpdateInterval / 2
					for trytime := 0; trytime < 2; trytime++ {
						err := j.S.UpdateService(nil)
						if err != nil {
							time.Sleep(waitSecond)
							continue
						} else {
							j.state.SetSuccess()
							return
						}
					}
					j.state.SetFail()
					log.Printf("[WARN]jobID:%v do updateservice fail, retry out of times", j.config.JobID)
					//					j.S.SetMachines(nil)
				}()

			} else if res == WARN {
				j.state.SetWarn()
				//log.Printf("[WARN]jobID:%v do health check warn", j.config.JobID)
			} else if res == FAIL {
				j.state.SetFail()
				//log.Printf("[WARN]jobID:%v do health check Fail", j.config.JobID)
			} else if res == INIT {
				err := j.S.OnlyUpdateService(nil)
				if err != nil {
					log.Printf("[INFO][INIT]jobID:%v domain not exist, only update called with nothing change.", j.config.JobID)
				} else {
					log.Printf("[INFO][INIT]jobID:%v domain exist, only update called change ttl to %v.", j.config.JobID, j.S.Ttl)
				}
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
				if res == WARN {
					log.Printf("[WARN]jobID:%v do health check warn", j.config.JobID)
				} else if res == FAIL {
					log.Printf("[WARN]jobID:%v do health check Fail", j.config.JobID)
				}
				j.state.LastCheckStatus = res
			}()
		}
	}
}
