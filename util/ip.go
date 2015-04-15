package util

import (
	"fmt"
	"net"
)

var privateBlocks []*net.IPNet

const ETCDPORT = "2379"
const ETCDMACHINES = "etcd.sdp"

//get ip by name ,use for etcd machines discoury
func GetIPByName(name string) []string {
	ns, err := net.LookupIP(name)
	if err != nil {
		fmt.Printf("no ips for %v\n", name)
		return nil
	} else {
		var ips []string
		for _, ip := range ns {
			ips = append(ips, ip.String())
		}
		return ips
	}
}

func init() {
	// Add each private block
	privateBlocks = make([]*net.IPNet, 3)
	_, block, err := net.ParseCIDR("10.0.0.0/8")
	if err != nil {
		panic(fmt.Sprintf("Bad cidr. Got %v", err))
	}
	privateBlocks[0] = block

	_, block, err = net.ParseCIDR("172.16.0.0/12")
	if err != nil {
		panic(fmt.Sprintf("Bad cidr. Got %v", err))
	}
	privateBlocks[1] = block

	_, block, err = net.ParseCIDR("192.168.0.0/16")
	if err != nil {
		panic(fmt.Sprintf("Bad cidr. Got %v", err))
	}
	privateBlocks[2] = block
}

func isPrivateIP(ip_str string) bool {
	ip := net.ParseIP(ip_str)
	for _, priv := range privateBlocks {
		if priv.Contains(ip) {
			return true
		}
	}
	return false
}

// GetPrivateIP is used to return the first private IP address
// associated with an interface on the machine
func GetPrivateIP() (net.IP, error) {
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		return nil, fmt.Errorf("Failed to get interface addresses: %v", err)
	}

	// Find private IPv4 address
	for _, rawAddr := range addresses {
		var ip net.IP
		switch addr := rawAddr.(type) {
		case *net.IPAddr:
			ip = addr.IP
		case *net.IPNet:
			ip = addr.IP
		default:
			continue
		}

		if ip.To4() == nil {
			continue
		}
		if !isPrivateIP(ip.String()) {
			continue
		}

		return ip, nil
	}

	return nil, fmt.Errorf("No private IP address found")
}
