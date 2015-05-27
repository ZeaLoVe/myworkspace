package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sdagent/sdagent"
	"sdagent/service"
	"sdagent/util"
	"strconv"
)

var begin int
var end int
var ttl int
var filepath string

func main() {
	flag.IntVar(&begin, "b", 8080, "port start number")
	flag.IntVar(&end, "e", 8139, "port stop number")
	flag.IntVar(&ttl, "t", 15, "default ttl")
	flag.StringVar(&filepath, "f", "/etc/sdconfig.json", "port stop number")
	flag.Parse()
	ip, _ := util.GetPrivateIP()
	agent := sdagent.SDAgent{}
	var sers []service.Service
	for i := begin; i <= end; i++ {
		ser := service.Service{}
		ser.Key = ip.String() + ":" + strconv.Itoa(i)
		ser.Port = uint64(i)
		ser.Ttl = uint64(ttl)
		var hc service.HealthCheck
		hc.HTTP = "http://" + ip.String() + ":" + strconv.Itoa(i)
		ser.Hc = append(ser.Hc, hc)
		sers = append(sers, ser)
	}
	agent.S = sers
	data, _ := json.Marshal(agent)
	//fmt.Printf(string(data))
	file, err := os.Create(filepath)
	if err != nil {
		fmt.Printf("get file error.\n")
	} else {
		file.Write(data)
	}
}
