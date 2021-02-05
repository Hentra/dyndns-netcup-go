package netcup

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

const (
	url = "https://ccp.netcup.net/run/webservice/servers/endpoint.php?JSON"
)

var (
	// ErrNoAPISessionid indicates that there is no session id available. This means that
	// you are probably not logged in.
	ErrNoAPISessionid = errors.New("netcup: There is no ApiSessionId. Are you logged in?")

	verbose = false
)

// Client represents a client to the netcup api.
type Client struct {
	client         *http.Client
	Customernumber int
	APIKey         string
	APIPassword    string
	APISessionid   string
}

// NewClient returns a new client by customernumber, apikey and apipassword
func NewClient(customernumber int, apikey, apipassword string) *Client {
	return &Client{
		Customernumber: customernumber,
		APIKey:         apikey,
		APIPassword:    apipassword,
		client:         http.DefaultClient,
	}
}

func (c *Client) do(req *Request) (*Response, error) {
	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Post(url, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	if !response.isSuccess() {
		return nil, errors.New(response.getFormattedError())
	}

	logInfo(response.getFormattedStatus())

	return &response, nil
}

// Login logs the client in the netcup api. This method should be issued before
// any other method.
func (c *Client) Login() error {
	var params = NewParams()
	params.AddParam("apikey", c.APIKey)
	params.AddParam("apipassword", c.APIPassword)
	params.AddParam("customernumber", strconv.Itoa(c.Customernumber))

	request := NewRequest("login", &params)

	response, err := c.do(request)
	if err != nil {
		return err
	}

	var loginResponse LoginResponse
	err = json.Unmarshal(response.ResponseData, &loginResponse)
	if err != nil {
		return err
	} else if loginResponse.APISessionid == "" {
		return errors.New("netcup: empty sessionid supplied")
	} else {
		c.APISessionid = loginResponse.APISessionid
	}

	return nil
}

// InfoDNSZone return the DNSZone for a specified domain
func (c *Client) InfoDNSZone(domainname string) (*DNSZone, error) {
	params, err := c.basicAuthParams(domainname)
	if err != nil {
		return nil, err
	}

	request := NewRequest("infoDnsZone", params)

	response, err := c.do(request)
	if err != nil {
		return nil, err
	}

	var dnsZone DNSZone
	err = json.Unmarshal(response.ResponseData, &dnsZone)
	if err != nil {
		return nil, err
	}

	return &dnsZone, nil
}

// InfoDNSRecords returns a DNSRecordSet for a specified domain
func (c *Client) InfoDNSRecords(domainname string) (*DNSRecordSet, error) {
	params, err := c.basicAuthParams(domainname)
	if err != nil {
		return nil, err
	}

	request := NewRequest("infoDnsRecords", params)

	response, err := c.do(request)
	if err != nil {
		return nil, err
	}

	var dnsRecordSet DNSRecordSet
	err = json.Unmarshal(response.ResponseData, &dnsRecordSet)
	if err != nil {
		return nil, err
	}

	return &dnsRecordSet, nil
}

// UpdateDNSZone updates the specified domain with a specified DNSZone
func (c *Client) UpdateDNSZone(domainname string, dnszone *DNSZone) error {
	params, err := c.basicAuthParams(domainname)
	if err != nil {
		return err
	}
	params.AddParam("dnszone", dnszone)
	request := NewRequest("updateDnsZone", params)

	_, err = c.do(request)
	if err != nil {
		return err
	}

	return nil
}

// UpdateDNSRecords updates the specified domain with a specified DNSRecordSet
func (c *Client) UpdateDNSRecords(domainname string, dnsRecordSet *DNSRecordSet) error {
	params, err := c.basicAuthParams(domainname)
	if err != nil {
		return err
	}

	params.AddParam("dnsrecordset", dnsRecordSet)
	request := NewRequest("updateDnsRecords", params)

	_, err = c.do(request)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) basicAuthParams(domainname string) (*Params, error) {
	if c.APISessionid == "" {
		return nil, ErrNoAPISessionid
	}

	params := NewParams()
	params.AddParam("apikey", c.APIKey)
	params.AddParam("apisessionid", c.APISessionid)
	params.AddParam("customernumber", strconv.Itoa(c.Customernumber))
	params.AddParam("domainname", domainname)

	return &params, nil
}

// SetVerbose sets the verboseness of the output. If set to true the response to every
// request will be send to stdout.
func SetVerbose(isVerbose bool) {
	verbose = isVerbose
}

func logInfo(msg string, v ...interface{}) {
	if verbose {
		log.Printf(msg, v...)
	}
}
