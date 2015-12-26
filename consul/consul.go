package consul

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/longnguyen11288/nomadic-router/Godeps/_workspace/src/github.com/hashicorp/consul/api"
	log "github.com/longnguyen11288/nomadic-router/Godeps/_workspace/src/gopkg.in/inconshreveable/log15.v2"
)

const MIN_PORT = 20000
const MAX_PORT = 60000

// Client to talk consul
type Client struct {
	endpoint     string
	consulClient *api.Client
}

//ServiceInfo to store array of services
type ServiceInfo struct {
	Name    string
	Address string
	ID      string
	Port    int
}

// NewClient Create new consul client
func NewClient(endpoint string) *Client {
	consulConfig := api.DefaultConfig()
	consulConfig.Address = endpoint
	consulClient, err := api.NewClient(consulConfig)
	if err != nil {
		log.Error("Error creating consul client", "error", err)
	}
	client := &Client{endpoint: endpoint, consulClient: consulClient}
	return client
}

// GetServices Get all services and return info
func (c *Client) GetServices() ([]*ServiceInfo, error) {
	serviceInfos := []*ServiceInfo{}
	services, _, err := c.consulClient.Catalog().Services(&api.QueryOptions{})
	if err != nil {
		return []*ServiceInfo{}, err
	}
	for service, tags := range services {
		for _, tag := range tags {
			consulService, _, _ := c.consulClient.Catalog().Service(service, tag, &api.QueryOptions{})
			for _, consulService := range consulService {
				serviceInfo := &ServiceInfo{
					Name:    consulService.ServiceName,
					Address: consulService.Address,
					ID:      consulService.ServiceID,
					Port:    consulService.ServicePort,
				}
				serviceInfos = append(serviceInfos, serviceInfo)
			}
		}
	}
	return serviceInfos, nil
}

func (c *Client) ClaimPort(port int) error {
	ports, _, err := c.consulClient.KV().Get("/nomadic/ports", &api.QueryOptions{})
	// Key doesn't exist return empty
	if ports == nil {
		return err
	}
	var nomadicPorts []int
	json.Unmarshal(ports.Value, &nomadicPorts)
	nomadicPorts = append(nomadicPorts, port)
	newPorts, err := json.Marshal(nomadicPorts)
	if err != nil {
		return err
	}
	ports.Value = newPorts
	c.consulClient.KV().Put(ports, &api.WriteOptions{})

	return nil
}

// GetPorts from consul
func (c *Client) GetPorts() []int {
	ports, _, err := c.consulClient.KV().Get("/nomadic/ports", &api.QueryOptions{})
	// Key doesn't exist return empty
	if ports == nil {
		kvp := &api.KVPair{Key: "/nomadic/ports", Value: []byte("[]")}
		c.consulClient.KV().Put(kvp, &api.WriteOptions{})
		return []int{}
	}
	var nomadicPorts []int
	err = json.Unmarshal(ports.Value, &nomadicPorts)
	if err != nil {
		fmt.Println(err)
	}
	return nomadicPorts
}

// GeneratePort used to generate random port
func (c *Client) GeneratePort() int {
	var randomNumber int
	usedPorts := c.GetPorts()
	randomNumber = random(MIN_PORT, MAX_PORT)
	for contains(usedPorts, randomNumber) {
		randomNumber = random(MIN_PORT, MAX_PORT)
	}
	return randomNumber
}

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
