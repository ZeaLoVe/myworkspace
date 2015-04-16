package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	. "myworkspace/util"
	"os"
	"strings"
	"sync"
	"time"
)

var ReqList []string

var QPS int
var ThreadNum int
var CheckRes bool
var RequestATread int
var CheckFile string
var CheckMap map[string]string //Domain Name -> list of address(ip or name)

var workers []Worker
var stopChan chan int

type Worker struct {
	ReqSendNum       int
	ResSuccessNum    int
	ResFailNum       int
	lenOfRequestList int
	interval         time.Duration
	lock             sync.Mutex
}

func (w *Worker) DoRequest() {
	w.ReqSendNum = w.ReqSendNum + 1
	//fmt.Printf("---%v %v %v", w.ReqSendNum, w.ResFailNum, w.ResSuccessNum)
	req := ReqList[w.ReqSendNum%w.lenOfRequestList]
	//fmt.Printf("sent......%v\n", req)
	ips := GetIPByName(req)
	if !CheckRes {
		if len(ips) != 0 {
			w.ResSuccessNum = w.ResSuccessNum + 1
		} else {
			w.ResFailNum = w.ResFailNum + 1
		}
	} else {
		for _, ip := range ips {
			if CheckMap[ip] == req {
				//fmt.Printf("CheckMap of %v:%v  req:%v\n", ip, CheckMap[ip], req)
				w.ResSuccessNum = w.ResSuccessNum + 1
				return
			}
		}
		w.ResFailNum = w.ResFailNum + 1
	}
}

func (w *Worker) Run() {
	for i := RequestATread; i > 0; i-- {
		time.Sleep(w.interval)
		w.DoRequest()
	}
	stopChan <- 10
}

func main() {
	flag.IntVar(&QPS, "q", 10000, "-q to set qps per thread")
	flag.IntVar(&ThreadNum, "t", 5, "-t to set worker thread,default set as 5")
	flag.IntVar(&RequestATread, "r", 10000, "-r to set the nums of request a worker will send ,default as 10000")
	flag.BoolVar(&CheckRes, "c", false, "-c to set whether check the result,default not")
	flag.StringVar(&CheckFile, "f", "", "-f to set check file, DNS(name ip) sets")
	flag.Parse()
	fmt.Printf("Set %v thread , do %v requests a thread, QPS:%v, check:%v ,checkFile:%v.", ThreadNum, RequestATread, QPS, CheckRes, CheckFile)
	if CheckFile == "" {
		fmt.Println("please use -f to give a checkfile")
		return
	} else {
		CheckMap = map[string]string{}
		file, err := os.Open(CheckFile)
		if err != nil {
			fmt.Println("open file error")
		}
		defer file.Close()
		buf := bufio.NewReader(file)
		for {
			line, err := buf.ReadString('\n')
			line = strings.Trim(line, "\r\n")
			if err == io.EOF {
				break
			}
			if err != nil && err != io.EOF {
				fmt.Println("read buf error")
				break

			} else {

				strs := strings.Split(line, " ")
				if len(strs) != 2 {
					fmt.Println("format error ,must be:service address")
					continue
				}
				CheckMap[strs[1]] = strs[0]        //make check map
				ReqList = append(ReqList, strs[0]) //make request list
			}
		}
		fmt.Println("Request List:")
		for _, str := range ReqList {
			fmt.Println(str)
		}
		fmt.Println("CheckMap:")
		for key, value := range CheckMap {
			fmt.Println(key + " " + value)
		}
	}
	stopChan = make(chan int, ThreadNum+1)
	for i := 0; i < ThreadNum; i++ {
		worker := new(Worker)
		workers = append(workers, *worker)
	}
	for i := 0; i < ThreadNum; i++ {
		workers[i].interval = time.Duration(float64(1.0/QPS)) * time.Second
		workers[i].lenOfRequestList = len(ReqList)
		workers[i].ReqSendNum = 0
		workers[i].ResFailNum = 0
		workers[i].ResSuccessNum = 0
		fmt.Printf("set thread internal:%v lenofRequest list:%v \n", workers[i].interval, workers[i].lenOfRequestList)
		go workers[i].Run()
	}
	//wait all worker stop
	ReqSendSum := 0
	ResFailSum := 0
	ResSuccessSum := 0
	for i, _ := range workers {
		<-stopChan
		//fmt.Printf("---%v %v %v", workers[i].ReqSendNum, workers[i].ResFailNum, workers[i].ResSuccessNum)
		ReqSendSum = ReqSendSum + workers[i].ReqSendNum
		ResFailSum = ResFailSum + workers[i].ResFailNum
		ResSuccessSum += workers[i].ResSuccessNum
	}
	fmt.Printf("Totally send %v request,Get %v, failed %v,success %v\n", ThreadNum*RequestATread, ResFailSum+ResSuccessSum, ResFailSum, ResSuccessSum)
	fmt.Printf("Usable rate: %v", float64(ResSuccessSum)/(float64(ThreadNum*RequestATread)))

}
