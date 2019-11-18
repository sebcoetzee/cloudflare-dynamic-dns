package cloudflare

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Zone struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type ListZonesResponse struct {
	APIResponse
	Result []Zone `json:"result"`
}

func (c *client) ListZones(zoneName string) (int, ListZonesResponse, error) {
	status, responseBody, err := c.httpClient.SendRequest(http.MethodGet, fmt.Sprintf("/zones?type=A&name=%s", zoneName), nil)
	if err != nil {
		return 500, ListZonesResponse{}, err
	}

	var zonesResponse ListZonesResponse
	err = json.Unmarshal(responseBody, &zonesResponse)
	if err != nil {
		return status, ListZonesResponse{}, err
	}

	return status, zonesResponse, err
}
