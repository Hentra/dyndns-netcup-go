package netcup

import (
    "log"
    "bytes"
    "encoding/json"
    "net/http"
    "io/ioutil"
    "strconv"
    "errors"
)

const (
    url = "https://ccp.netcup.net/run/webservice/servers/endpoint.php?JSON"
)

var (
    ErrNoApiSessionid = errors.New("netcup: There is no ApiSessionId. Are you logged in?")
)

type Client struct {
    client *http.Client
    Customernumber int
    ApiKey string
    ApiPassword string
    ApiSessionid string
}


func NewClient(customernumber int, apikey, apipassword string) *Client {
    return &Client{
        Customernumber: customernumber,
        ApiKey: apikey,
        ApiPassword: apipassword,
        client: http.DefaultClient,
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

    log.Println(response.getFormattedStatus())

    return &response, nil
}

func (c *Client) Login() error {
    var params = NewParams()
    params.AddParam("apikey", c.ApiKey)
    params.AddParam("apipassword", c.ApiPassword)
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
    } else if loginResponse.ApiSessionid == "" {
        return errors.New("netcup: empty sessionid supplied")
    } else {
        c.ApiSessionid = loginResponse.ApiSessionid
    }


    return nil
}

func (c *Client) InfoDnsZone(domainname string) (error, *DNSZone) {
    params := c.basicAuthParams(domainname)
    request := NewRequest("infoDnsZone", params)

    response, err := c.do(request)
    if err != nil {
        return err, nil
    }

    var dnsZone DNSZone
    err = json.Unmarshal(response.ResponseData, &dnsZone)
    if err != nil {
        return err, nil
    }

    return nil, &dnsZone
}

func (c *Client) InfoDnsRecords(domainname string) (error, *DNSRecordSet) {
    params := c.basicAuthParams(domainname)
    request := NewRequest("infoDnsRecords", params)

    response, err := c.do(request)
    if err != nil {
        return err, nil
    }

    var dnsRecordSet DNSRecordSet
    err = json.Unmarshal(response.ResponseData, &dnsRecordSet)
    if err != nil {
        return err, nil
    }

    return nil, &dnsRecordSet
}

func (c *Client) UpdateDnsZone(domainname string, dnszone *DNSZone) error {
    params := c.basicAuthParams(domainname)
    params.AddParam("dnszone", dnszone)
    request := NewRequest("updateDnsZone", params)

    _, err := c.do(request)
    if err != nil {
        return err
    }

    return nil
}

func (c *Client) UpdateDnsRecords(domainname string, dnsRecordSet *DNSRecordSet) error {
    params := c.basicAuthParams(domainname)
    params.AddParam("dnsrecordset", dnsRecordSet)
    request := NewRequest("updateDnsRecords", params)

    _, err := c.do(request)
    if err != nil {
        return err
    }

    return nil
}

func (c *Client) basicAuthParams(domainname string) *Params {
    params := NewParams()
    params.AddParam("apikey", c.ApiKey)
    params.AddParam("apisessionid", c.ApiSessionid)
    params.AddParam("customernumber", strconv.Itoa(c.Customernumber))
    params.AddParam("domainname", domainname)

    return &params
}

