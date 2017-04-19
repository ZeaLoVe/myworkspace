/*
SDAgent versiong 0.1
Service config file needed use -f=filapath  default sdconfig.json
Etcd machines get from DNS use -d=DNS       default zealove.xyz
Etcd port set              use -p=port      default 2379
Etcd protocol set          use -h=protocol  default http://
Reload interval set        use -t=num       default 10 unit minute
*/

package main

import (
	"flag"
	"log"
	"os"
	. "sdagent/backends"
	. "sdagent/sdagent"
	. "sdagent/util"
	"time"
)

func env(key, def string) string {
	if x := os.Getenv(key); x != "" {
		return x
	}
	return def
}

//1.3.3 如果health check 结果为INIT未执行，则在5s后再次调用域名更新的操作
//1.3.4 修复域名注册的时候存在大写字母的问题
//1.3.5 增加health check 的默认超时时间从3秒到10秒
//1.4.1 增加etcd账号
const Version = "1.4.1"

//ETCDPORT\ETCDDOMAIN\MODIFYINTERVAL come from backends
var CONFIGFILE string
var PIDFILEPATH string
var SHOWVERSION bool

func main() {
	flag.StringVar(&CONFIGFILE, "f", env("SDAGENT_CONFIGFILE", "sdconfig.json"), "Path of config file")
	flag.StringVar(&ETCDDOMAIN, "d", env("SDAGENT_ETCDDOMAIN", "etcd.sdp"), "Name for DNS request of etcd machines")
	flag.StringVar(&ETCDPROTOCOL, "h", env("SDAGENT_ETCDPROTOCOL", "http://"), "etcd client protocol")
	flag.StringVar(&ETCDPORT, "p", env("SDAGENT_ETCDPORT", "2379"), "etcd client port")
	flag.StringVar(&PIDFILEPATH, "m", "", "gen pid file ,use for monit")
	flag.StringVar(&ETCDACCOUNT, "u", env("ETCD_ACCOUNT", ""), "etcd account")
	flag.StringVar(&ETCDPASSWORD, "w", env("ETCD_PASSWORD", ""), "etcd password")
	flag.IntVar(&MODIFYINTERVAL, "t", 1, "Reload Check Interval")
	flag.BoolVar(&SHOWVERSION, "v", false, "sdagent version")
	flag.Parse()

	if SHOWVERSION {
		log.Println("sdagent version: ", Version)
		return
	}

	if err := GenPidFile(PIDFILEPATH, "sdagent.pid"); err != nil {
		log.Printf("[WARN]Gen pid file error with %v.\n", err)
	}

	if CONFIGFILE != "" {
		log.Printf("[INFO]SDAgent use file:%v  for configure.\n", CONFIGFILE)
		agent := NewAgent(CONFIGFILE)
		if agent == nil {
			log.Printf("[WARN]Can't init from given config file:%v .Nothing run.\n", CONFIGFILE)
		} else {

			//Band etcd machine with backend，reload will not change it，if changed need restart
			DefaultBackend.SetMachines(nil)
			agent.Start()
		}
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

		for { //main thread sleep
			time.Sleep(time.Duration(MODIFYINTERVAL) * time.Minute)
		}
	}
}
