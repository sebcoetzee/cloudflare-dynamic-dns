package cloudflare

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type DNSRecord struct {
	ID      string `json:"id,omitempty"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Content string `json:"content"`
	Proxied bool   `json:"proxied"`
}

type ListDNSRecordsResponse struct {
	APIResponse
	Result []DNSRecord `json:"result"`
}

func (c *client) ListDNSRecords(zoneID string, name string) (int, ListDNSRecordsResponse, error) {
	status, responseBody, err := c.httpClient.SendRequest(http.MethodGet, fmt.Sprintf("/zones/%s/dns_records?name=%s", zoneID, name), nil)
	if err != nil {
		return 500, ListDNSRecordsResponse{}, err
	}

	var dnsRecordsResponse ListDNSRecordsResponse
	err = json.Unmarshal(responseBody, &dnsRecordsResponse)
	if err != nil {
		return status, ListDNSRecordsResponse{}, err
	}

	return status, dnsRecordsResponse, err
}

func (c *client) CreateDNSRecord(zoneID string, dnsRecord DNSRecord) (int, APIResponse, error) {
	bodyBytes, err := json.Marshal(&dnsRecord)
	if err != nil {
		return 500, APIResponse{}, err
	}

	status, responseBody, err := c.httpClient.SendRequest(http.MethodPost, fmt.Sprintf("/zones/%s/dns_records", zoneID), bytes.NewReader(bodyBytes))
	if err != nil {
		return 500, APIResponse{}, err
	}

	var result APIResponse
	err = json.Unmarshal(responseBody, &result)
	if err != nil {
		return 500, APIResponse{}, err
	}

	return status, result, err
}

func (c *client) UpdateDNSRecord(zoneID string, dnsRecord DNSRecord) (int, APIResponse, error) {
	bodyBytes, err := json.Marshal(&dnsRecord)
	if err != nil {
		return 500, APIResponse{}, err
	}

	status, responseBody, err := c.httpClient.SendRequest(http.MethodPut, fmt.Sprintf("/zones/%s/dns_records/%s", zoneID, dnsRecord.ID), bytes.NewReader(bodyBytes))
	if err != nil {
		return 500, APIResponse{}, err
	}

	var result APIResponse
	err = json.Unmarshal(responseBody, &result)
	if err != nil {
		return 500, APIResponse{}, err
	}

	return status, result, err
}
