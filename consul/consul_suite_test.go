package consul_test

import (
	"testing"

	. "github.com/longnguyen11288/nomadic-router/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "github.com/longnguyen11288/nomadic-router/Godeps/_workspace/src/github.com/onsi/gomega"
)

func TestConsul(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Consul Test Suite")
}
