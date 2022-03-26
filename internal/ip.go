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
func GetAddrInfo(ipv4 bool, ipv6 bool) (*AddrInfo, error) {
	adresses := &AddrInfo{}

	if ipv4 {
		address, err := getIPv4()
		if err != nil {
			return nil, err
		}
		adresses.IPv4 = address
	}

	if ipv6 {
		address, err := getIPv6()
		if err != nil {
			return nil, err
		}

		adresses.IPv6 = address
	}

	return adresses, nil
}

func getIPv4() (string, error) {
	return do("https://api.ipify.org?format=text")
}

func getIPv6() (string, error) {
	return do("https://api6.ipify.org?format=text")
}

func do(url string) (string, error) {
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
