package avengersclient

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Default avengers-backend URL
const HostURL string = "http://localhost:8000"

type Client struct {
	HostURL    string
	HTTPClient *http.Client
}

func NewClient(host *string) (*Client, error) {
	c := Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		HostURL:    HostURL,
	}
	if host != nil {
		c.HostURL = *host
	}
	return &c, nil
}

//DoRequest will actually make a call to avengers-backend
func (c *Client) DoRequest(req *http.Request) ([]byte, error) {
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}
