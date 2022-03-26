package netcup

import (
	"encoding/json"
)

// Response represents a response from the netcup api.
type Response struct {
	ServerRequestID string          `json:"serverrequestid"`
	ClientRequestID string          `json:"clientrequestid"`
	Action          string          `json:"action"`
	Status          string          `json:"status"`
	StatusCode      int             `json:"statuscode"`
	ShortMessage    string          `json:"shortmessage"`
	LongMessage     string          `json:"longmessage"`
	ResponseData    json.RawMessage `json:"responsedata"`
}

// LoginResponse represents the response from the netcup api which is unique
// to the LoginRequest.
type LoginResponse struct {
	APISessionid string `json:"apisessionid"`
}

// DNSZone represents a dns zone.
type DNSZone struct {
	DomainName   string `json:"name"`
	TTL          string `json:"ttl"`
	Serial       string `json:"serial"`
	Refresh      string `json:"refresh"`
	Retry        string `json:"retry"`
	Expire       string `json:"expire"`
	DNSSecStatus bool   `json:"dnssecstatus"`
}

// DNSRecord represents a dns record.
type DNSRecord struct {
	ID           string `json:"id"`
	Hostname     string `json:"hostname"`
	Type         string `json:"type"`
	Priority     string `json:"priority"`
	Destination  string `json:"destination"`
	DeleteRecord bool   `json:"deleterecord"`
	State        string `json:"state"`
}

// DNSRecordSet represents a dns record set.
type DNSRecordSet struct {
	DNSRecords []DNSRecord `json:"dnsrecords"`
}

// NewDNSRecordSet return a new DNSRecordSet containing specified DNSRecords.
func NewDNSRecordSet(records []DNSRecord) *DNSRecordSet {
	return &DNSRecordSet{
		DNSRecords: records,
	}
}

// GetRecordOccurences returns the amount of times a specified hostname with a dnstype
// occures in the DNSRecordSet.
func (r *DNSRecordSet) GetRecordOccurences(hostname, dnstype string) int {
	result := 0
	for _, record := range r.DNSRecords {
		if record.Hostname == hostname && record.Type == dnstype {
			result++
		}
	}
	return result
}

// GetRecord returns DNSRecord that matches both the name and dnstype specified or nil
// if its not inside the DNSRecordSet.
func (r *DNSRecordSet) GetRecord(name, dnstype string) *DNSRecord {
	for _, record := range r.DNSRecords {
		if record.Hostname == name && record.Type == dnstype {
			return &record
		}
	}
	return nil
}

// NewDNSRecord returns a new DNSRecord with specified hostname, dnstype and destination.
func NewDNSRecord(hostname, dnstype, destination string) *DNSRecord {
	return &DNSRecord{
		Hostname:    hostname,
		Type:        dnstype,
		Destination: destination,
	}
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
