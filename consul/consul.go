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

// PortMapping struct
type PortMapping struct {
	Name string `json:"name"`
	Port int    `json:"port"`
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

func (c *Client) ClaimPort(name string, port int) error {
	ports, _, err := c.consulClient.KV().Get("nomadic/ports", &api.QueryOptions{})
	// Key doesn't exist return empty
	if ports == nil {
		return err
	}
	var nomadicPorts []PortMapping
	json.Unmarshal(ports.Value, &nomadicPorts)
	serviceMapping := PortMapping{Name: name, Port: port}
	nomadicPorts = append(nomadicPorts, serviceMapping)
	newPorts, err := json.Marshal(nomadicPorts)
	if err != nil {
		return err
	}
	ports.Value = newPorts
	c.consulClient.KV().Put(ports, &api.WriteOptions{})

	return nil
}

// GetPorts from consul
func (c *Client) GetPorts() []PortMapping {
	ports, _, err := c.consulClient.KV().Get("nomadic/ports", &api.QueryOptions{})
	// Key doesn't exist return empty
	if ports == nil {
		kvp := &api.KVPair{Key: "nomadic/ports", Value: []byte("[]")}
		_, err := c.consulClient.KV().Put(kvp, &api.WriteOptions{})
		if err != nil {
			fmt.Println(err)
		}
		return []PortMapping{}
	}
	var nomadicPorts []PortMapping
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

func (c *Client) GeneratePortOrGetCurrent(name string) (int, bool) {
	ports := c.GetPorts()
	for _, mapping := range ports {
		if mapping.Name == name {
			return mapping.Port, false
		}
	}
	return c.GeneratePort(), true
}

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

func contains(s []PortMapping, e int) bool {
	for _, a := range s {
		if a.Port == e {
			return true
		}
	}
	return false
}

