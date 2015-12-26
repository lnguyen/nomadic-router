package main

import (
	"flag"
	"time"

	"github.com/longnguyen11288/nomadic-router/core"
)

var (
	consul = flag.String("c", "127.0.0.1:8500", "URL for Consul HTTP API")
	gorb   = flag.String("g", "http://127.0.0.1:4672", "URL for Consul HTTP API")
	ip     = flag.String("ip", "", "IP to map to")
)

func main() {
	flag.Parse()
	for {
		core.PopulateIVS(*gorb, *consul, *ip)
		time.Sleep(time.Second * 15)
	}
}
