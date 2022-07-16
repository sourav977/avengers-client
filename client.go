package avengersclient

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"gopkg.in/retry.v1"
)

const (
	HostURL              = "http://localhost:8000" // Default avengers-backend URL
	DEFAULT_HTTP_TIMEOUT = 10 * time.Second
	InitialBackoffDelay  = 1 * time.Second
	MaxBackoffDelay      = 16 * time.Second // MaxBackoffDelay is a DelayType which increases delay between consecutive retries
	MaxRetries           = 5
)

type Client struct {
	HostURL    string
	HTTPClient *http.Client
	clock      retry.Clock
	strategy   retry.Strategy
}

var Strategy = retry.LimitCount(MaxRetries, retry.Exponential{
	Initial:  InitialBackoffDelay,
	Factor:   2,
	MaxDelay: MaxBackoffDelay,
})

func NewClient(host *string, clock retry.Clock, strategy retry.Strategy) (*Client, error) {
	c := Client{
		HTTPClient: &http.Client{Timeout: DEFAULT_HTTP_TIMEOUT},
		HostURL:    HostURL,
	}
	if host != nil {
		c.HostURL = *host
	}
	if clock != nil {
		c.clock = clock
	} else {
		c.clock = nil
	}
	if strategy != nil {
		c.strategy = strategy
	} else {
		c.strategy = Strategy
	}
	return &c, nil
}

//DoRequest will actually make a call to avengers-backend
func (c *Client) DoRequest(req *http.Request) ([]byte, error) {
	res, err := c.doRetryableHTTPRequest(req)
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

func (c *Client) doRetryableHTTPRequest(req *http.Request) (resp *http.Response, err error) {
	for a := retry.Start(Strategy, nil); a.Next(); {
		resp, err = c.HTTPClient.Do(req)
		if err != nil {
			return
		}
		defer resp.Body.Close()

		if a.More() && resp.StatusCode/100 == 5 {
			fmt.Printf("Attemtping to retry request, to %s", resp.Request.URL)
			continue
		}
		return
	}
	return nil, fmt.Errorf("Retry strategy did not allow any attempts")
}
