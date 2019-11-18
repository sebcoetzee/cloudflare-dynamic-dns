package cloudflare

import (
	"github.com/sebcoetzee/cloudflare-dynamic-dns/config"
	"github.com/sebcoetzee/cloudflare-dynamic-dns/httpclient"
)

func NewClient(config config.Config) Client {
	return &client{
		httpClient: httpclient.NewClient().
			BaseURL("https://api.cloudflare.com/client/v4").
			BearerToken(config.APIToken),
	}
}

type Client interface {
	ListZones(zoneName string) (int, ListZonesResponse, error)
	ListDNSRecords(zoneID string, name string) (int, ListDNSRecordsResponse, error)
	CreateDNSRecord(zoneID string, dnsRecord DNSRecord) (int, APIResponse, error)
	UpdateDNSRecord(zoneID string, dnsRecord DNSRecord) (int, APIResponse, error)
}

type client struct {
	httpClient httpclient.Client
}
