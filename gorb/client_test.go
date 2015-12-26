package gorb_test

import (
	. "github.com/longnguyen11288/nomadic-router/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "github.com/longnguyen11288/nomadic-router/Godeps/_workspace/src/github.com/onsi/gomega"
	"github.com/longnguyen11288/nomadic-router/Godeps/_workspace/src/github.com/onsi/gomega/ghttp"
	. "github.com/longnguyen11288/nomadic-router/gorb"
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

	Describe("creating a new service ", func() {
		BeforeEach(func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("PUT", "/service/foo"),
				),
			)
		})

		It("create a new service", func() {
			svc := Service{
				Host:       "10.244.234.64",
				Port:       445,
				Protocol:   "tcp",
				Method:     "rr",
				Persistent: true,
			}
			client.NewService("foo", svc)
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
	})

	Describe("creating a new backend", func() {
		BeforeEach(func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("PUT", "/service/foo/bar"),
				),
			)
		})

		It("create a new backend", func() {
			backend := Backend{
				Host:   "10.244.234.64",
				Port:   8500,
				Weight: 50,
				Method: "nat",
			}
			client.NewBackend("bar", "foo", backend)
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
	})

	Describe("get a service", func() {
		BeforeEach(func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/service/foo"),
				),
			)
		})

		It("get a service", func() {
			client.GetService("foo")
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
	})

	Describe("get a backend", func() {
		BeforeEach(func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/service/foo/bar"),
				),
			)
		})

		It("get a backend", func() {
			client.GetBackend("bar", "foo")
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
	})

})
