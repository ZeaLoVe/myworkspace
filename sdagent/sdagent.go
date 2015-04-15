package sdagent

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	. "myworkspace/service"
	"time"
)

const Version = "0.1"
const Usage = ` SDAgent versiong 0.1
Service config file needed in the same folder,named sdconfig.json
Config file need be json formation.
`

type SDAgent struct {
	S    []Service `json:"services,omitempty"`
	Jobs []Job     `json:"-"`

	stopAgentChan chan uint64
}

var Agent SDAgent

func init() {
	Agent = SDAgent{}
	Agent.stopAgentChan = make(chan uint64, 128)
	Agent.LoadConfig("sdconfig.json")
	for i, _ := range Agent.S {
		job := NewJob()
		job.S = Agent.S[i]
		job.S.SetDefault()
		job.S.InitService()
		Agent.Jobs = append(Agent.Jobs, *job)
	}
}

func (sda *SDAgent) LoadConfig(filename string) error {
	config, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal("[ERR]Can't load config file,start agent error!\n")
		return errors.New("Can't load config file,start agent error!")
	}
	if err := json.Unmarshal(config, sda); err != nil {
		log.Fatal("[ERR]Unmarsh to JSON fail,start agent error!\n")
		return errors.New("Unmarsh to JSON fail,start agent error!")
	}
	return nil
}

func (sda *SDAgent) Start() {
	countRun := 0
	countFail := 0
	for i, _ := range sda.Jobs {
		if !sda.Jobs[i].CanRun() {
			log.Printf("[DEBUG]jobID:%v miss something, will not run.\n", sda.Jobs[i].config.JobID)
			countFail++
		} else {
			sda.Jobs[i].SetConfig()
			go sda.Jobs[i].Run()
			countRun++
		}
	}
	log.Printf("[DEBUG]Totally %v Jobs in config, run %v jobs, %v failed.\n", countRun+countFail, countRun, countFail)
}

func (sda *SDAgent) StopAll() {
	log.Println("[DEBUG]Agent Stop all begin!")
	for i, _ := range sda.Jobs {
		if sda.Jobs[i].config.JOBSTATE == RUNNING {
			sda.stopAgentChan <- STOPCHANNUM
			sda.Jobs[i].stopChan <- STOPCHANNUM //stop job
			<-sda.Jobs[i].stopChan              //wait for stop
		}
	}
	//close(sda.stopAgentChan) pannic here....
	log.Println("[DEBUG]Agent all job stopped!")
}

func (sda *SDAgent) StopJob(job *Job) {
	if job.config.JOBSTATE == RUNNING {
		job.stopChan <- STOPCHANNUM
		<-job.stopChan
		log.Println("[DEBUG]A job stopped!")
	} else {
		log.Println("[DEBUG]StopJob cmd sent to a job not RUNNING")
	}
}

func (sda *SDAgent) StartJob(job *Job) {
	if !job.CanRun() {
		log.Printf("[DEBUG]jobID:%v miss something, will not run\n", job.config.JobID)
	} else {
		if job.config.JOBSTATE == RUNNING {
			log.Printf("[DEBUG]jobID:%v already Running, can't start again\n", job.config.JobID)
		}
		if job.config.JOBSTATE == PREPARE {
			job.SetConfig()
			go job.Run()
		}
	}
}

func (sda *SDAgent) AutoCheck(i int) {
	timeout := time.After(sda.Jobs[i].config.UpdateInterval)
	//timeout := time.After(sda.Jobs[i].config.UpdateInterval)
	heartbeat := time.Tick(sda.Jobs[i].config.UpdateInterval / 2)
	for {
		select {
		case <-timeout:
			log.Printf("[DEBUG]jobID:%v timeout will restart.\n", sda.Jobs[i].config.JobID)
			sda.StartJob(&sda.Jobs[i])
		case <-heartbeat:
			if keep, ok := <-sda.Jobs[i].keepAliveChan; ok {
				if keep == KEEPALIVENUM {
					timeout = time.After(sda.Jobs[i].config.UpdateInterval) //reflesh timeout
					log.Printf("[DEBUG]Agent get jobID:%v heartbeat\n", sda.Jobs[i].config.JobID)
				}
			}
		case <-sda.stopAgentChan:
			log.Printf("[DEBUG]Agent check jobID:%v stop\n", sda.Jobs[i].config.JobID)
			return
		}
	}
}

func (sda *SDAgent) Run() {
	for i, _ := range sda.Jobs {
		if sda.Jobs[i].CanRun() {
			go sda.AutoCheck(i)
		}
	}
}
