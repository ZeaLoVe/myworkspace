package sdagent

import (
	. "myworkspace/service"
	"testing"
	"time"
)

func TestJobSetConfig(t *testing.T) {
	var ser Service
	ser.SetDefault()
	var job Job
	job.s = &ser
	job.SetConfig()
	if job.config.UpdateInterval != time.Duration(job.s.Ttl/2)*time.Second {
		t.Fatalf("job Test SetConfig UpdateInterval fail")
	}
	if job.config.JobID != job.s.Key {
		t.Fatalf("job Test SetConfig JobID fail")
	}
}

func TestJobRun(t *testing.T) {
	testjob := new(Job)
	ser := new(Service)
	ser.SetDefault()
	testjob.s = ser
	testjob.SetConfig()
	go testjob.Run()

	time.Sleep(30 * time.Second)
	testjob.stopChan <- 10
	time.Sleep(10 * time.Second)

}
