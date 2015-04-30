package sdagent

import (
	"log"
	. "myworkspace/service"
	"testing"
	"time"
)

func TestJobState(t *testing.T) {
	j := new(Job)
	if j.config.JOBSTATE == PREPARE {
		log.Printf("Job State init right")
	} else {
		log.Printf("Job State init error")
	}
	j.SetJobState(RUNNING)
	if j.config.JOBSTATE == RUNNING {
		log.Printf("Job State Set RUNNING right")
	} else {
		log.Printf("Job State Set RUNNING error")
	}
	j.SetJobState(PREPARE)
	if j.config.JOBSTATE == PREPARE {
		log.Printf("Job State Set PREPARE right")
	} else {
		log.Printf("Job State Set PREPARE error")
	}
}

func TestJobSetConfig(t *testing.T) {
	var ser Service
	var job Job
	ser.SetDefault()
	ser.InitService()
	job.S = ser
	job.SetConfig()
	if job.config.JOBSTATE != READY {
		log.Printf("job serconfig error test fail")
	}
	job.SetConfig()
	if job.config.JOBSTATE != READY {
		log.Printf("job serconfig fail")
	}
	if job.config.UpdateInterval != time.Duration(job.S.Ttl-1)*time.Second {
		t.Fatalf("job Test SetConfig UpdateInterval fail")
	}
	//log.Printf("%v ----- %v", job.config.JobID, job.S.Key)
	if job.config.JobID != job.S.Key {
		t.Fatalf("job Test SetConfig JobID fail")
	}
}

func TestJobRun(t *testing.T) {
	testjob := Job{}
	var ser Service
	ser.SetDefault()
	ser.InitService()
	//ser.Dump()
	testjob.S = ser
	testjob.SetConfig()
	go testjob.Run()

	time.Sleep(15 * time.Second)
	testjob.stopChan <- 10
	time.Sleep(1 * time.Second)
	log.Println("TestJobRun stop")
}
