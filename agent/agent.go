package main

import (
	. "myworkspace/sdagent"
	"time"
)

func main() {
	defer Agent.StopAll()
	Agent.Start()
	Agent.Run()
	time.Sleep(10 * time.Second)
}
