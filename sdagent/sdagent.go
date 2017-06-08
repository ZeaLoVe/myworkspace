package sdagent

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	. "sdagent/service"
	"sdagent/util"
	"time"
)

type SDAgent struct {
	S    []Service `json:"services,omitempty"`
	Jobs []*Job    `json:"-"`
}

func NewAgent(config string) *SDAgent {
	agent := new(SDAgent)
	err := agent.LoadConfig(config)
	if err != nil {
		log.Printf("[ERR]Can't load Config File with err:%v\n", err)
		return nil
	} else {
		for i, _ := range agent.S {
			job := NewJob()
			for j, _ := range agent.S[i].Hc {
				agent.S[i].Hc[j].SetDefault()
			}
			job.S = agent.S[i]
			job.S.SetDefault()
			job.S.InitService()
			agent.Jobs = append(agent.Jobs, job)
		}
	}
	return agent
}

func (sda *SDAgent) LoadConfig(filename string) error {
	config, err := ioutil.ReadFile(filename)
	if err != nil {
		return errors.New("Can't load config file,Can't start agent!")
	}
	if err := json.Unmarshal(config, sda); err != nil {
		return errors.New("Parse JSON file fail, Cat't start agent!")
	}
	return nil
}

//restart agent's jobs by modify config file
func (sda *SDAgent) Reload(filename string) (*SDAgent, error) {
	ok, err := util.CheckModify(filename)
	if err != nil {
		return nil, err
	} else {
		if !ok {
			return nil, fmt.Errorf("File not changed recently")
		} else {
			tmp := NewAgent(filename)
			if tmp == nil {
				return nil, fmt.Errorf("Start new agent error while reload")
			} else {
				tmp.Start()
				if sda != nil {
					sda.StopAll()
				}
				return tmp, nil
			}
		}
	}
}

func (sda *SDAgent) Start() {
	countRun := 0
	countFail := 0
	for i, _ := range sda.Jobs {
		if !sda.Jobs[i].CanRun() {
			log.Printf("[ERR]jobID:%v config has something wrong, will not run.\n", sda.Jobs[i].config.JobID)
			countFail++
		} else {
			sda.Jobs[i].SetConfig()
			go sda.Jobs[i].Run()
			countRun++
		}
	}
	log.Printf("[INFO]Totally %v Jobs in config, run %v jobs, %v failed.\n", countRun+countFail, countRun, countFail)
}

func (sda *SDAgent) StopAll() {
	errCh := make(chan error, 2)
	go func() {
		time.Sleep(5 * time.Second) //set default timeout
		errCh <- fmt.Errorf("%Stop All timeout")
	}()
	go func() {
		count := 0
		for i, _ := range sda.Jobs {
			if sda.Jobs[i].config.JOBSTATE == RUNNING {
				sda.Jobs[i].stopChan <- STOPCHANNUM //stop job
				//sda.Jobs[i].S.OnlyUpdateService(nil) //only update
				<-sda.Jobs[i].stopChan //wait for stop
				count++
			}
		}
		errCh <- nil
		log.Printf("[INFO]Agent stopped %v jobs!\n", count)
	}()

	err := <-errCh
	if err != nil {
		log.Println("[ERR]Agent Stop All timeout")
	}
}

func (sda *SDAgent) StopJob(job *Job) {
	if job.config.JOBSTATE == RUNNING {
		job.stopChan <- STOPCHANNUM
		<-job.stopChan
	} else {
		log.Println("[WARN]StopJob cmd sent to a job not RUNNING")
	}
}

func (sda *SDAgent) StartJob(job *Job) {
	if !job.CanRun() {
		log.Printf("[WARN]jobID:%v miss something, will not run\n", job.config.JobID)
	} else {
		if job.config.JOBSTATE == RUNNING {
			log.Printf("[WARN]jobID:%v already running, can't start again\n", job.config.JobID)
		}
		if job.config.JOBSTATE == PREPARE {
			job.SetConfig()
			go job.Run()
		}
	}
}
