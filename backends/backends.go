package backends

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"path"
	. "sdagent/util"
	"strings"
	"time"

	etcdClient "github.com/coreos/etcd/client"
	"golang.org/x/net/context"
)

//sdagent default transport setting
var SDagentTransport etcdClient.CancelableTransport = &http.Transport{
	Proxy: http.ProxyFromEnvironment,
	Dial: (&net.Dialer{
		Timeout:   5 * time.Second,
		KeepAlive: 5 * time.Second,
	}).Dial,
	TLSHandshakeTimeout: 5 * time.Second,
}

type Backend struct {
	client etcdClient.KeysAPI
}

var DefaultBackend = Backend{}

func init() {
	DefaultBackend.SetMachines(nil)
}

//Get key write to etcd
func GenKey(name string) string {
	tmpList := strings.Split(name, ".")
	for i, j := 0, len(tmpList)-1; i < j; i, j = i+1, j-1 {
		tmpList[i], tmpList[j] = tmpList[j], tmpList[i]
	}
	return path.Join(append([]string{"/skydns/"}, tmpList...)...)
}

func (backend *Backend) SetMachines(newMachine []string) error {
	var tmpMachines []string
	if len(newMachine) == 0 || (len(newMachine) == 1 && newMachine[0] == "") {
		tmpMachines = GetIPByName(ETCDDOMAIN)
		if tmpMachines == nil {
			return fmt.Errorf("DNS can't got any etcd machines")
		}
		for i, machine := range tmpMachines {
			tmpMachines[i] = ETCDPROTOCOL + machine + ":" + ETCDPORT
		}
	} else { //replace
		tmpMachines = newMachine
	}
	log.Println("etcd machines set:", tmpMachines)
	cfg := etcdClient.Config{
		Endpoints: tmpMachines,
		Transport: SDagentTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}
	c, err := etcdClient.New(cfg)
	if err != nil {
		return err
	}
	backend.client = etcdClient.NewKeysAPI(c)
	return nil
}

func (backend *Backend) OnlyUpdate(key string, value string, ttl uint64) error {
	if backend.client == nil {
		if err := backend.SetMachines(nil); err != nil {
			return err
		}
	}
	errCh := make(chan error, 2)
	go func() {
		var errSet error
		if _, errSet = backend.client.Set(context.Background(),
			key,
			value,
			&etcdClient.SetOptions{PrevExist: etcdClient.PrevExist, TTL: time.Duration(ttl) * time.Second}); errSet == nil {
			//fmt.Println("update success")
		}
		errCh <- errSet
	}()

	go func() {
		time.Sleep(time.Duration(3) * time.Second)
		errCh <- fmt.Errorf("etcd only update timeout")
	}()

	err := <-errCh
	return err
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
		if _, errSet = backend.client.Set(context.Background(),
			key,
			value,
			&etcdClient.SetOptions{TTL: time.Duration(ttl) * time.Second}); errSet == nil {
			//fmt.Println("set success")
		}
		errCh <- errSet
	}()

	go func() {
		time.Sleep(time.Duration(3) * time.Second)
		errCh <- fmt.Errorf("etcd timeout")
	}()

	err := <-errCh
	return err
}

func (backend *Backend) CheckKey(key string) (bool, error) {
	if backend.client == nil {
		if err := backend.SetMachines(nil); err != nil {
			return false, err
		}
	}
	resp, err := backend.client.Get(context.Background(), key, nil)
	if err != nil {
		if strings.Contains(err.Error(), "Key not found") {
			return false, err
		}
	} else {
		if resp.Node.TTL >= 0 {
			return true, nil
		} else {
			return false, fmt.Errorf("key is out of ttl")
		}
	}
	return true, nil
}
