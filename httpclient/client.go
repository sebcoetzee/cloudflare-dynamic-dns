package httpclient

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func NewClient() Client {
	return &client{
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

type Client interface {
	BaseURL(baseURL string) Client
	BearerToken(token string) Client
	SendRequest(method string, url string, body io.Reader) (int, []byte, error)
}

type client struct {
	bearerToken string
	baseURL     string
	httpClient  *http.Client
}

func (c *client) BaseURL(baseURL string) Client {
	c.baseURL = baseURL
	return c
}

func (c *client) BearerToken(token string) Client {
	c.bearerToken = token
	return c
}

func (c *client) SendRequest(method string, url string, body io.Reader) (int, []byte, error) {
	url = fmt.Sprintf("%s%s", c.baseURL, url)

	logEntry := logrus.WithField("url", url)

	if body != nil {
		bodyBytes, err := ioutil.ReadAll(body)
		if err != nil {
			return 500, nil, err
		}
		body = bytes.NewReader(bodyBytes)
		logEntry = logEntry.WithField("request", string(bodyBytes))
	}

	logEntry.Debug("starting request")
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return 500, nil, err
	}

	if c.bearerToken != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.bearerToken))
	}

	timeBefore := time.Now()
	res, err := c.httpClient.Do(req)
	if err != nil {
		return 500, nil, err
	}
	totalTime := time.Since(timeBefore)

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 500, nil, err
	}

	logEntry.WithFields(logrus.Fields{
		"status":       res.StatusCode,
		"response":     string(responseBody),
		"request_time": totalTime,
	}).Debug("finished request")

	return res.StatusCode, responseBody, nil
}
