package sdagent

import (
	"fmt"
	"log"
	"testing"
	"time"
)

var Agent SDAgent

func init() {
	Agent = SDAgent{}
	Agent.LoadConfig("sdconfig.json")
	for i, _ := range Agent.S {
		job := NewJob()
		job.S = Agent.S[i]
		job.S.SetDefault()
		job.S.InitService()
		Agent.Jobs = append(Agent.Jobs, job)
	}
}

func TestNewAgent(t *testing.T) {
	if agent := NewAgent("sdconfig.json"); agent != nil {
		log.Println("New Agent call success")
	}
}

func TestStart(t *testing.T) {
	fmt.Println("test Agent Start ")
	Agent.Start()
	time.Sleep(5 * time.Second)
	if Agent.Jobs[0].config.JOBSTATE != RUNNING {
		log.Printf("Job Not run")
	}
	if Agent.Jobs[1].config.JOBSTATE != RUNNING {
		log.Printf("Job Not run")
	}
}

func TestRunAndStop(t *testing.T) {
	fmt.Println("test RunAndStop start ")
	Agent.StartJob(Agent.Jobs[0])
	time.Sleep(5 * time.Second)
	Agent.StopJob(Agent.Jobs[0])
	time.Sleep(5 * time.Second)
	if Agent.Jobs[0].config.JOBSTATE != PREPARE {
		log.Printf("Stop dont change jobstate")
	}
	Agent.StartJob(Agent.Jobs[1])
	Agent.StartJob(Agent.Jobs[0])
	time.Sleep(5 * time.Second)
	Agent.StopAll()
}

func TestAgentRun(t *testing.T) {
	Agent.Start()
	time.Sleep(5 * time.Second)
	Agent.StopJob(Agent.Jobs[0])
	time.Sleep(5 * time.Second)
	Agent.StopJob(Agent.Jobs[1])
	time.Sleep(5 * time.Second)
}
