package consulstore

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/consul/api"
	"github.com/longnguyen11288/nomadic-router/store"
)

type ConsulStore struct {
	client *api.Client
}

func New(address string) *ConsulStore {
	config := api.DefaultConfig()
	config.Address = address
	client, _ := api.NewClient(config)
	return &ConsulStore{
		client: client,
	}
}

//GetService return a service from consul
func (c *ConsulStore) GetService(serviceId string) (*store.Service, error) {
	name, err := c.getValue("/services/" + serviceId + "/name")
	if err != nil {
		return nil, err
	}
	host, err := c.getValue("/services/" + serviceId + "/host")
	if err != nil {
		return nil, err
	}
	portString, err := c.getValue("/services/" + serviceId + "/port")
	if err != nil {
		return nil, err
	}
	protocol, err := c.getValue("/services/" + serviceId + "/protocol")
	if err != nil {
		return nil, err
	}
	scheduler, err := c.getValue("/services/" + serviceId + "/scheduler")
	if err != nil {
		return nil, err
	}
	port, err := strconv.ParseInt(string(portString), 10, 16)
	if err != nil {
		return nil, err
	}
	var destinations []store.Destination
	destPathList, _, err := c.client.KV().Keys("/services/"+serviceId+"/destinations/", "/", nil)
	for _, destPath := range destPathList {
		host, err := c.getValue(destPath + "host")
		if err != nil {
			return nil, err
		}
		portString, err := c.getValue(destPath + "port")
		if err != nil {
			return nil, err
		}
		weightString, err := c.getValue(destPath + "weight")
		if err != nil {
			return nil, err
		}
		mode, err := c.getValue(destPath + "mode")
		if err != nil {
			return nil, err
		}
		port, err := strconv.ParseInt(string(portString), 10, 16)
		if err != nil {
			return nil, err
		}
		weight, err := strconv.ParseInt(string(weightString), 10, 32)
		if err != nil {
			return nil, err
		}
		destination := store.Destination{
			Host:   string(host),
			Port:   uint16(port),
			Weight: int32(weight),
			Mode:   string(mode),
		}
		destinations = append(destinations, destination)
	}
	return &store.Service{
		Name:         string(name),
		Host:         string(host),
		Port:         uint16(port),
		Protocol:     string(protocol),
		Scheduler:    string(scheduler),
		Destinations: destinations,
	}, nil
}

// GetServicesName returns all service names
func (c *ConsulStore) GetServicesName() ([]string, error) {
	names, _, err := c.client.KV().Keys("/services/", "/", nil)
	if err != nil {
		return nil, err
	}
	for i, k := range names {
		k = strings.TrimSuffix(k, "/")
		names[i] = strings.TrimPrefix(k, "services/")
	}
	return names, nil
}

//UpsertService adding service to consul
func (c *ConsulStore) UpsertService(svc store.Service) error {
	err := c.setValue(svc.Path()+"/name", []byte(svc.Name))
	if err != nil {
		return err
	}
	err = c.setValue(svc.Path()+"/host", []byte(svc.Host))
	if err != nil {
		return err
	}
	err = c.setValue(svc.Path()+"/port", []byte(strconv.Itoa(int(svc.Port))))
	if err != nil {
		return err
	}
	err = c.setValue(svc.Path()+"/protocol", []byte(svc.Protocol))
	if err != nil {
		return err
	}
	err = c.setValue(svc.Path()+"/scheduler", []byte(svc.Scheduler))
	if err != nil {
		return err
	}
	return nil
}

// DeleteService remove service
func (c *ConsulStore) DeleteService(svc store.Service) error {
	_, err := c.client.KV().DeleteTree(svc.Path(), nil)
	if err != nil {
		return err
	}
	return nil
}

//func (c *ConsulStore) GetDestinations(svc store.Service) (*[]Destination, error) {
//}

// UpsertDestination add destination if not exist
func (c *ConsulStore) UpsertDestination(svc store.Service, dst store.Destination) error {
	path := svc.Path() + dst.Path()
	err := c.setValue(path+"/host", []byte(dst.Host))
	if err != nil {
		return err
	}
	err = c.setValue(path+"/port", []byte(strconv.Itoa(int(dst.Port))))
	if err != nil {
		return err
	}
	err = c.setValue(path+"/weight", []byte(strconv.Itoa(int(dst.Weight))))
	if err != nil {
		return err
	}
	err = c.setValue(path+"/mode", []byte(dst.Mode))
	if err != nil {
		return err
	}
	return nil
}

// DeleteDestination remove a destination
func (c *ConsulStore) DeleteDestination(svc store.Service, dst store.Destination) error {
	_, err := c.client.KV().DeleteTree(svc.Path()+dst.Path(), nil)
	if err != nil {
		return err
	}
	return nil
}

// Subscribe to watch for changes
func (c *ConsulStore) Subscribe(changes chan interface{}) error {
	return nil
}

// Flush consul store
func (c *ConsulStore) Flush() error {
	return nil
}

func (c *ConsulStore) getValue(path string) ([]byte, error) {
	kvpair, _, err := c.client.KV().Get(path, nil)
	return kvpair.Value, err
}

func (c *ConsulStore) setValue(key string, value []byte) error {
	kp := &api.KVPair{
		Key:   key,
		Value: value,
	}
	time, err := c.client.KV().Put(kp, nil)
	if err != nil {
		log.Printf("Consul [%d]: %s", time, err)
		return fmt.Errorf("Consul [%d]: %s", time, err)
	}
	return nil
}
