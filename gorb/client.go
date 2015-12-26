package gorb

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

// Client to talk to gorb
type Client struct {
	endpoint   string
	httpclient *http.Client
}

// Service struct to new service info
type Service struct {
	Host       string `json:"host"`
	Port       uint16 `json:"port"`
	Protocol   string `json:"protocol"`
	Method     string `json:"method"`
	Persistent bool   `json:"persistent"`
}

// Backend struct to hold backend service info
type Backend struct {
	Host   string `json:"host"`
	Port   uint16 `json:"port"`
	Weight int32  `json:"weight"`
	Method string `json:"method"`
}

// NewClient create new gorb client
func NewClient(endpoint string) *Client {
	c := &Client{endpoint: endpoint, httpclient: http.DefaultClient}
	return c
}

// NewService create a new service
func (c *Client) NewService(name string, svc Service) error {
	url := c.endpoint + "/service/" + name
	body, err := json.Marshal(svc)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PUT", url, strings.NewReader(string(body)))
	if err != nil {
		return err
	}
	res, err := c.httpclient.Do(req)
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return err
	}
	return nil
}

// NewBackend create a new service
func (c *Client) NewBackend(name string, serviceName string, backend Backend) error {
	url := c.endpoint + "/service/" + serviceName + "/" + name
	body, err := json.Marshal(backend)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PUT", url, strings.NewReader(string(body)))
	if err != nil {
		return err
	}
	res, err := c.httpclient.Do(req)
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return err
	}
	return nil
}

// GetService get a service, check to see if service is there
func (c *Client) GetService(name string) error {
	url := c.endpoint + "/service/" + name
	req, err := http.NewRequest("GET", url, strings.NewReader(""))
	if err != nil {
		return err
	}
	res, err := c.httpclient.Do(req)
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return errors.New("Unable to find service")
	}
	return nil
}

// GetBackend get a backend, check to see if service is there
func (c *Client) GetBackend(name string, svcName string) error {
	url := c.endpoint + "/service/" + svcName + "/" + name
	req, err := http.NewRequest("GET", url, strings.NewReader(""))
	if err != nil {
		return err
	}
	res, err := c.httpclient.Do(req)
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return errors.New("Unable to find backend")
	}
	return nil
}
