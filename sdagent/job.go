package sdagent

import ()

type Job struct {
	r        *Register
	s        *Service
	stopChan chan uint64 //check whether job is down
}

func (j *Job) Run() {

}
