package consul_test

import (
	. "github.com/longnguyen11288/nomadic-router/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "github.com/longnguyen11288/nomadic-router/consul"
	//. "github.com/onsi/gomega"
	"github.com/longnguyen11288/nomadic-router/Godeps/_workspace/src/github.com/onsi/gomega/ghttp"
)

var _ = Describe("Client API", func() {

	var server *ghttp.Server
	var client *Client

	BeforeEach(func() {
		server = ghttp.NewServer()
		client = NewClient(server.URL())
	})

	AfterEach(func() {
		//shut down the server between tests
		server.Close()
	})

})
