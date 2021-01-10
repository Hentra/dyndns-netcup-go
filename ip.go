package main

import (
	"io/ioutil"
	"net/http"
)

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
