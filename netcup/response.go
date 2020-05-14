package netcup

import (
    "encoding/json"
)

type Response struct {
    ServerRequestID string `json:"serverrequestid"`
    ClientRequestID string `json:"clientrequestid"`
    Action string `json:"action"`
    Status string `json:"status"`
    StatusCode int `json:"statuscode"`
    ShortMessage string `json:"shortmessage"`
    LongMessage string `json:"longmessage"`
    ResponseData json.RawMessage `json:"responsedata"`
}

type LoginResponse struct {
    ApiSessionid string `json:"apisessionid"`
}

type DNSZone struct {
    DomainName string `json:"name"`
    TTL string `json:"ttl"`
    Serial string `json:"serial"`
    Refresh string `json:"refresh"`
    Retry string `json:"retry"`
    Expire string `json:"expire"`
    DNSSecStatus bool `json:"dnssecstatus"`
}

type DNSRecord struct {
    ID string `json:"id"`
    Hostname string `json:"hostname"`
    Type string `json:"type"`
    Priority string `json:"priority"`
    Destination string `json:"destination"`
    DeleteRecord bool `json:"deleterecord"`
    State string `json:"state"`
}

type DNSRecordSet struct {
    DNSRecords []DNSRecord `json:"dnsrecords"`
}

func NewDNSRecordSet(records []DNSRecord) *DNSRecordSet {
    return &DNSRecordSet{
        DNSRecords: records,
    }
}

func (r *DNSRecordSet) GetARecordOccurences(hostname string) int {
    result := 0
    for _, record := range r.DNSRecords {
        if record.Hostname == hostname && record.Type == "A" {
            result++
        }
    }
    return result
}

func (r *DNSRecordSet) GetARecord(name string) (*DNSRecord, bool) {
    for _, record := range r.DNSRecords {
        if record.Hostname == name && record.Type == "A" {
            return &record, true
        }
    }
    return nil, false
}

func (r *Response) isError() bool {
    return r.Status == "error"
}

func (r *Response) isSuccess() bool {
    return r.Status == "success"
}

func (r *Response) getFormattedError() string {
    return "netcup: " + r.ShortMessage + " Reason: " + r.LongMessage
}

func (r *Response) getFormattedStatus() string {
    return "netcup: [" + r.Status + "] " + r.LongMessage
}
