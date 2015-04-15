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
		Agent.Jobs = append(Agent.Jobs, *job)
	}
}

func (sda *SDAgent) LoadConfig(filename string) error {
	config, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal("Can't load config file,start agent error!")
		return errors.New("Can't load config file,start agent error!")
	}
	if err := json.Unmarshal(config, sda); err != nil {
		log.Fatal("Unmarsh to JSON fail,start agent error!")
		return errors.New("Unmarsh to JSON fail,start agent error!")
	}
	return nil
}

func (sda *SDAgent) Start() {
	countRun := 0
	countFail := 0
	for i, _ := range sda.Jobs {
		if !sda.Jobs[i].CanRun() {
			log.Printf("jobID:%v miss something, will not run\n", sda.Jobs[i].config.JobID)
			countFail++
		} else {
			sda.Jobs[i].SetConfig()
			go sda.Jobs[i].Run()
			countRun++
		}
	}
	log.Printf("Totally %v Jobs in config, run %v jobs, %v failed.\n", countRun+countFail, countRun, countFail)
}

func (sda *SDAgent) StopAll() {
	sda.StopAgent()
	log.Println("Agent Stop all begin!")
	for i, _ := range sda.Jobs {
		if sda.Jobs[i].config.JOBSTATE == RUNNING {
			sda.Jobs[i].stopChan <- 10 //stop job
			<-sda.Jobs[i].stopChan     //wait for stop
		}
	}
	close(sda.stopAgentChan)
	log.Println("Agent all job stopped!")
}

func (sda *SDAgent) StopJob(job *Job) {
	if job.config.JOBSTATE == RUNNING {
		job.stopChan <- 10
		log.Println("Agent Stop a job cmd sent!")
		<-job.stopChan
		log.Println("Agent a job stopped!")
	} else {
		log.Println("StopJob cmd sent to a job not RUNNING")
	}
}

func (sda *SDAgent) StartJob(job *Job) {
	if !job.CanRun() {
		log.Printf("jobID:%v miss something, will not run\n", job.config.JobID)
	} else {
		if job.config.JOBSTATE == RUNNING {
			log.Printf("jobID:%v already Running, can't start again\n", job.config.JobID)
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
			log.Printf("jobID:%v timeout will restart.\n", sda.Jobs[i].config.JobID)
			sda.StartJob(&sda.Jobs[i])
		case <-heartbeat:
			if keep, ok := <-sda.Jobs[i].keepAliveChan; ok {
				if keep == KEEPALIVENUM {
					timeout = time.After(sda.Jobs[i].config.UpdateInterval) //reflesh timeout
					log.Printf("Agent get jobID:%v heartbeat\n", sda.Jobs[i].config.JobID)
				}
			}
		case <-sda.stopAgentChan:
			log.Printf("Agent check jobID:%v stop\n", sda.Jobs[i].config.JobID)
			return
		}
	}
}

func (sda *SDAgent) StopAgent() {
	for i := len(sda.S); i > 0; i-- {
		sda.stopAgentChan <- 34
		sda.stopAgentChan <- 34
	}
}

func (sda *SDAgent) Run() {
	for i, _ := range sda.Jobs {
		if sda.Jobs[i].CanRun() {
			go sda.AutoCheck(i)
		}
	}
}
