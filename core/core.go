package core

import (
	"github.com/longnguyen11288/nomadic-router/consul"
	"github.com/longnguyen11288/nomadic-router/gorb"
)

// PopulateIVS set up IVS tables
func PopulateIVS(gorbEndpoint, consulEndpoint, IP string) error {
	gorbClient := gorb.NewClient(gorbEndpoint)
	consulClient := consul.NewClient(consulEndpoint)
	services, err := consulClient.GetServices()
	if err != nil {
		return err
	}
	for _, service := range services {
		err := gorbClient.GetService(service.Name)
		if err != nil {
			port := consulClient.GeneratePortOrGetCurrent(service.Name)
			svc := gorb.Service{
				Host:       IP,
				Port:       uint16(port),
				Protocol:   "tcp",
				Method:     "rr",
				Persistent: true,
			}
			err = gorbClient.NewService(service.Name, svc)
			if err != nil {
				return err
			}
			err = consulClient.ClaimPort(service.Name, port)
			if err != nil {
				return err
			}
		}
		err = gorbClient.GetBackend(service.ID, service.Name)
		if err != nil {
			backend := gorb.Backend{
				Host:   service.Address,
				Port:   uint16(service.Port),
				Weight: 50,
				Method: "nat",
			}
			err = gorbClient.NewBackend(service.ID, service.Name, backend)
			if err != nil {
				return nil
			}
		}
	}

	return nil
}

