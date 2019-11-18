package ipfy

import (
	"net"
	"net/http"

	"github.com/sebcoetzee/cloudflare-dynamic-dns/httpclient"
)

func NewClient() Client {
	return &client{
		apiv4Client: httpclient.NewClient().
			BaseURL("https://api.ipify.org"),
		apiv6Client: httpclient.NewClient().
			BaseURL("https://api6.ipify.org"),
	}
}

type Client interface {
	GetIPv4() (int, net.IP, error)
	GetIPv6() (int, net.IP, error)
}

type client struct {
	apiv4Client httpclient.Client
	apiv6Client httpclient.Client
}

func (c *client) GetIPv4() (int, net.IP, error) {
	status, responseBody, err := c.apiv4Client.SendRequest(http.MethodGet, "", nil)
	if err != nil {
		return 500, nil, err
	}

	if status != 200 {
		return status, responseBody, nil
	}

	return status, net.ParseIP(string(responseBody)), nil
}

func (c *client) GetIPv6() (int, net.IP, error) {
	status, responseBody, err := c.apiv6Client.SendRequest(http.MethodGet, "", nil)
	if err != nil {
		return 500, nil, err
	}

	if status != 200 {
		return status, responseBody, nil
	}

	return status, net.ParseIP(string(responseBody)), nil
}
