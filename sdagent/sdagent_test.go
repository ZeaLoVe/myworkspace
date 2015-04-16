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

func TestNewAgent(t *testing.T) {
	if agent := NewAgent("sdconfig.json"); agent != nil {
		log.Println("New Agent call success")
	}
}

func TestAgentLoadConfig(t *testing.T) {
	fmt.Println("test Agent load config ")
	//for _, s := range Agent.S {
	//	s.Dump()
	//}
	if len(Agent.Jobs) == 0 {
		fmt.Println("no jobs")
	}
	//for _, j := range Agent.Jobs {
	//	j.S.Dump()
	//}
}

func TestStart(t *testing.T) {
	fmt.Println("test Agent Start ")
	Agent.Start()
	time.Sleep(5 * time.Second)
	if Agent.Jobs[0].config.JOBSTATE != PREPARE {
		log.Printf("Stop ALL dont change jobstate")
	}
	if Agent.Jobs[1].config.JOBSTATE != PREPARE {
		log.Printf("Stop ALL dont change jobstate")
	}
}

func TestRunAndStop(t *testing.T) {
	fmt.Println("test RunAndStop start ")
	Agent.StartJob(&Agent.Jobs[0])
	time.Sleep(5 * time.Second)
	Agent.StopJob(&Agent.Jobs[0])
	time.Sleep(5 * time.Second)
	if Agent.Jobs[0].config.JOBSTATE != PREPARE {
		log.Printf("Stop dont change jobstate")
	}
	Agent.StartJob(&Agent.Jobs[1])
	Agent.StartJob(&Agent.Jobs[0])
	time.Sleep(10 * time.Second)
	Agent.StopAll()
}

func TestAgentRun(t *testing.T) {
	Agent.Start()
	Agent.Run()
	time.Sleep(10 * time.Second)
	Agent.StopJob(&Agent.Jobs[0])
	time.Sleep(5 * time.Second)
	Agent.StopJob(&Agent.Jobs[1])
	time.Sleep(10 * time.Second)
}
