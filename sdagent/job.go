package sdagent

import ()

type JobConfig struct {
	UpdateInterval uint64
	CheckInterval  uint64
}

type Job struct {
	r             *Register
	s             *Service
	keepAliveChan chan uint64 //check whether job is down
	stopChan      chan uint64 //get cmd to stop
	config        JobConfig
}

func (j *Job) Run() {

}
