package backends

import (
	"fmt"
	"github.com/coreos/go-etcd/etcd"
	"path"
	. "sdagent/util"
	"strings"
	"time"
)

type Backend struct {
	timeout int
	client  *etcd.Client
}

//Get key write to etcd
func GenKey(name string) string {
	tmpList := strings.Split(name, ".")
	for i, j := 0, len(tmpList)-1; i < j; i, j = i+1, j-1 {
		tmpList[i], tmpList[j] = tmpList[j], tmpList[i]
	}

	return path.Join(append([]string{"/skydns/"}, tmpList...)...)
}

func (backend *Backend) SetTimeout(timeout int) {
	backend.timeout = timeout
}

func (backend *Backend) SetMachines(newMachine []string) error {
	if backend.timeout == 0 {
		backend.timeout = 5
	}
	if len(newMachine) == 0 || (len(newMachine) == 1 && newMachine[0] == "") {
		tmpMachines := GetIPByName(ETCDDOMAIN)
		if tmpMachines == nil {
			return fmt.Errorf("DNS can't got any etcd machines")
		}
		for i, machine := range tmpMachines {
			tmpMachines[i] = ETCDPROTOCOL + machine + ":" + ETCDPORT
		}
		backend.client = etcd.NewClient(tmpMachines)
	} else { //replace
		backend.client = etcd.NewClient(newMachine)
	}
	return nil
}

func (backend *Backend) UpdateKV(key string, value string, ttl uint64) error {
	if backend.client == nil {
		if err := backend.SetMachines(nil); err != nil {
			return err
		}
	}
	errCh := make(chan error, 2)
	go func() {
		var errSet error
		if _, errSet = backend.client.Update(key, value, ttl); errSet == nil {
			//update success
			//fmt.Println("update success")
		} else {
			//fmt.Println("update fail do set")
			if _, errSet = backend.client.Set(key, value, ttl); errSet == nil {
				//set success
				//fmt.Println("Set success")
			}
		}
		errCh <- errSet
	}()

	go func() {
		time.Sleep(time.Duration(backend.timeout) * time.Second)
		errCh <- fmt.Errorf("etcd timeout")
	}()

	err := <-errCh
	return err
}
