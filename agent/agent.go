package main

import (
	"flag"
	"fmt"
	"log"
	. "myworkspace/sdagent"
	"strings"
	"time"
)

const Version = "0.1"
const Usage = ` SDAgent versiong 0.1
Service config file needed use -f=filapath.default sdconfig.json
Etcd Machines address needed use -e="http://ip:port".default "http://192.168.181.16:2379"
`

var ConfigFile string
var EtcdMachines string

func main() {
	flag.StringVar(&ConfigFile, "f", "sdconfig.json", "path of config file")
	flag.StringVar(&EtcdMachines, "e", "http://192.168.181.16:2379", "etcd address")
	flag.Parse()

	if ConfigFile != "" {
		log.Printf("Agent will use file:%v  for configure.\nEtcd machines: %v to setup.\n", ConfigFile, EtcdMachines)
		agent := NewAgent(ConfigFile)
		if agent == nil {
			fmt.Printf("Can't init from given config file:%v .Check the config file to make it right.\n", ConfigFile)
		} else {
			tmpEtcd := strings.Split(EtcdMachines, ",")
			for i, _ := range agent.Jobs {
				if agent.Jobs[i].S.Machines == "" {
					//better add regex check "http://ip:port"
					if len(tmpEtcd) == 1 && tmpEtcd[0] == "" {
						agent.Jobs[i].S.SetMachines(nil)
					} else {
						agent.Jobs[i].S.SetMachines(tmpEtcd)
					}
				}
			}
			defer agent.StopAll()
			agent.Start()
			agent.Run()
			for { // can run a http server here
				time.Sleep(time.Hour * 1)
			}
		}
	}
}
