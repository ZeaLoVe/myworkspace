// httpdemo project main.go
package main

import (
	"bufio"
	"fmt"
	"github.com/coreos/go-etcd/etcd"
	"io"
	"net"
	"os"
	"path"
	"strconv"
	"strings"
)

var machines []string //etcd machines
var usage string

func getipByName(name string) []net.IP {
	ns, err := net.LookupIP(name)
	if err != nil {
		fmt.Println("no ips for the name")
		return ns
	} else {
		fmt.Println("get ips for " + name)
		return ns
	}
}

func getData(key string) string {
	client := etcd.NewClient(machines)
	resp, err := client.Get(key, false, true)
	if err != nil {
		return err.Error()
	} else {
		return "key:" + resp.Node.Key + "value:" + resp.Node.Value + "expiration:" + resp.Node.Expiration.String()
	}
}

func setData(key string, value string, ttl uint64) error {
	client := etcd.NewClient(machines)
	_, err := client.Set(key, value, ttl)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func getService(name string) string {
	tmpList := strings.Split(name, ".")
	for i, j := 0, len(tmpList)-1; i < j; i, j = i+1, j-1 {
		tmpList[i], tmpList[j] = tmpList[j], tmpList[i]
	}
	key := path.Join(append([]string{"/skydns/"}, tmpList...)...)
	return getData(key)
}

func setService(name string, ip string, ttl uint64) error {
	tmpList := strings.Split(name, ".")
	for i, j := 0, len(tmpList)-1; i < j; i, j = i+1, j-1 {
		tmpList[i], tmpList[j] = tmpList[j], tmpList[i]
	}
	key := path.Join(append([]string{"/skydns/"}, tmpList...)...)
	value := "{\"host\":\"" + ip + "\"}"
	fmt.Println("insert key: " + key)
	fmt.Println("insert value: " + value)
	err := setData(key, value, ttl)
	return err
}

func main() {
	machines = []string{"http://192.168.181.16:2379"} //set default
	usage = `
	SKYDNSTOOL version 1.0
	Usage:
	set name ip ttl
	forexample:set x1.mongo.sd.sdp 192.168.1.1 3600
	get name
	forexample:get x1.mongo.sd.sdp 
	batch set by file
	forexample:file RR.txt 
	RR.txt--------------->
	example.com 192.168.10.111 3600
	test.sdp.cn 192.198.1.11 3600
	x1.mongo.cn 192.198.1.13 3600
	
	other command has no results.
	`
	fmt.Printf(usage)
	for {

		fmt.Print("SKYDNSTOOL>")

		cmdReader := bufio.NewReader(os.Stdin)
		if cmdStr, err := cmdReader.ReadString('\n'); err == nil {

			cmdStr = strings.Trim(cmdStr, "\r\n")
			str := strings.Split(cmdStr, " ")

			switch str[0] {
			case "set":
				if len(str) != 4 {
					fmt.Println(usage)
					break
				}
				if ttl, err := strconv.Atoi(str[3]); err != nil {
					fmt.Println(err.Error())
					break
				} else {
					err := setService(str[1], str[2], uint64(ttl))
					if err == nil {
						fmt.Println("Set Record Succes")
					}
				}
			case "get":
				if len(str) != 2 {
					fmt.Println(usage)
					break
				}
				if result := getService(str[1]); result != "" {
					fmt.Println(result)
				}
			case "file":
				var count int
				count = 0
				file, err := os.Open(str[1])
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
						if len(strs) != 3 {
							fmt.Println("format error ,must be:key value ttl")
							continue
						}
						ttl, err := strconv.Atoi(strs[2])
						if err == nil {
							err := setService(strs[0], strs[1], uint64(ttl))
							if err == nil {
								count = count + 1
							}
						}
					}
				}
				fmt.Printf("%v lines of records set success.\n", count)
			default:
				fmt.Println(usage)
			}
		}
	}

}
