package internal

import (
	"io/ioutil"
	"net/http"
)

// AddrInfo represents the ip addresses of the host
type AddrInfo struct {
	IPv4 string
	IPv6 string
}

// GetAddrInfo retrieves an AddrInfo instance
func GetAddrInfo(ipv4 bool, ipv6 bool, ipv4service string, ipv6service string) (*AddrInfo, error) {
	adresses := &AddrInfo{}

	if ipv4 {
		address, err := getIP(ipv4service)
		if err != nil {
			return nil, err
		}
		adresses.IPv4 = address
	}

	if ipv6 {
		address, err := getIP(ipv6service)
		if err != nil {
			return nil, err
		}

		adresses.IPv6 = address
	}

	return adresses, nil
}

func getIP(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(ip), nil
}
