/*
SDAgent versiong 0.1
Service config file needed use -f=filapath  default sdconfig.json
Etcd machines get from DNS use -d=DNS       default zealove.xyz
Etcd port set              use -p=port      default 2379
Reload interval set        use -t=num       default 30 unit minute
*/

package main

import (
	"flag"
	"log"
	. "myworkspace/sdagent"
	. "myworkspace/util"
	"os"
	"time"
)

func env(key, def string) string {
	if x := os.Getenv(key); x != "" {
		return x
	}
	return def
}

const Version = "0.1"

//ETCDPORT\ETCDDOMAIN\MODIFYINTERVAL come from util
var CONFIGFILE string

func main() {

	flag.StringVar(&CONFIGFILE, "f", env("SDAGENT_CONFIGFILE", "sdconfig.json"), "Path of config file")
	flag.StringVar(&ETCDDOMAIN, "d", env("SDAGENT_ETCDDOMAIN", "zealove.xyz"), "Name for DNS request of etcd")
	flag.StringVar(&ETCDPORT, "p", env("SDAGENT_ETCDPORT", "2379"), "etcd client port")
	flag.IntVar(&MODIFYINTERVAL, "t", 30, "Reload Check Interval")
	flag.Parse()
	if CONFIGFILE != "" {
		log.Printf("[INFO]SDAgent use file:%v  for configure.\n", CONFIGFILE)
		agent := NewAgent(CONFIGFILE)
		if agent == nil {
			log.Printf("[ERR]Can't init from given config file:%v .Check the config file to make it right.\n", CONFIGFILE)
		} else {

			agent.Start()

			go func() {
				for {
					time.Sleep(time.Duration(MODIFYINTERVAL) * time.Minute)
					tmp, err := agent.Reload(CONFIGFILE)
					if err == nil && tmp != nil {
						agent = tmp
						log.Println("[RELOAD]Reload success.")
					} else {
						//log.Println("[RELOAD]Rugular Checked Config File, Not Reload.")
					}
				}
			}()

			for {
				time.Sleep(time.Duration(MODIFYINTERVAL) * time.Minute)
				CONFIGFILE = env("SDAGENT_CONFIGFILE", CONFIGFILE)
				ETCDDOMAIN = env("SDAGENT_ETCDDOMAIN", ETCDDOMAIN)
				ETCDPORT = env("SDAGENT_ETCDPORT", ETCDPORT)
			}
		}
	}
}
