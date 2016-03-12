package store

import (
	"fmt"
	"net"
	"syscall"

	"github.com/longnguyen11288/nomadic-router/ipvs"
)

const (
	NatMode    = ipvs.DFForwardMasq
	TunnelMode = ipvs.DFForwardTunnel
	RouteMode  = ipvs.DFForwardRoute
)

type Service struct {
	Name         string `valid:"alphanum,required"`
	Host         string `valid:"required"`
	Port         uint16 `valid:"required"`
	Protocol     string `valid:"required"`
	Scheduler    string `valid:"required"`
	Destinations []Destination
}

type Destination struct {
	Host   string `valid:"required"`
	Port   uint16 `valid:"required"`
	Weight int32
	Mode   string `valid:"required"`
}

// Path of service
func (svc Service) Path() string {
	return fmt.Sprintf("/services/%s", svc.GetID())
}

// GetID return id of service
func (svc Service) GetID() string {
	return fmt.Sprintf("%v-%v-%v", svc.Host, svc.Port, svc.Protocol)
}

// Path of destination
func (dst Destination) Path() string {
	return fmt.Sprintf("/destinations/%s", dst.GetID())
}

// GetID return id of destination
func (dst Destination) GetID() string {
	return fmt.Sprintf("%v-%v", dst.Host, dst.Port)
}

func stringToIPProto(s string) ipvs.IPProto {
	var value ipvs.IPProto
	if s == "udp" {
		value = syscall.IPPROTO_UDP
	} else {
		value = syscall.IPPROTO_TCP
	}

	return value
}

//MarshalJSON ...
func ipProtoToString(proto ipvs.IPProto) string {
	var value string

	if proto == syscall.IPPROTO_UDP {
		value = "udp"
	} else {
		value = "tcp"
	}

	return value
}

func stringToDestinationFlags(s string) ipvs.DestinationFlags {
	var flag ipvs.DestinationFlags

	switch s {
	case "nat":
		flag = NatMode
	case "tunnel":
		flag = TunnelMode
	default:
		// Default is Direct Routing
		flag = RouteMode
	}

	return flag
}

//MarshalJSON ...
func destinationFlagsToString(flags ipvs.DestinationFlags) string {
	var value string

	switch flags {
	case NatMode:
		value = "nat"
		// *flags =
	case TunnelMode:
		value = "tunnel"
	default:
		// Default is Direct Routing
		value = "route"
	}

	return value
}

func (s Service) ToIpvsService() ipvs.Service {
	destinations := []*ipvs.Destination{}

	for _, dst := range s.Destinations {
		destinations = append(destinations, dst.ToIpvsDestination())
	}

	return ipvs.Service{
		Address:      net.ParseIP(s.Host),
		Port:         s.Port,
		Protocol:     stringToIPProto(s.Protocol),
		Scheduler:    s.Scheduler,
		Destinations: destinations,
	}
}

func (d Destination) ToIpvsDestination() *ipvs.Destination {
	return &ipvs.Destination{
		Address: net.ParseIP(d.Host),
		Port:    d.Port,
		Weight:  d.Weight,
		Flags:   stringToDestinationFlags(d.Mode),
	}
}

func NewServiceRequest(s *ipvs.Service) Service {
	destinations := []Destination{}

	for _, dst := range s.Destinations {
		destinations = append(destinations, newDestinationRequest(dst))
	}

	return Service{
		Host:         s.Address.String(),
		Port:         s.Port,
		Protocol:     ipProtoToString(s.Protocol),
		Scheduler:    s.Scheduler,
		Destinations: destinations,
	}
}

func newDestinationRequest(d *ipvs.Destination) Destination {
	return Destination{
		Host:   d.Address.String(),
		Port:   d.Port,
		Weight: d.Weight,
		Mode:   destinationFlagsToString(d.Flags),
	}
}
