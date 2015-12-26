package gorb_test

import (
	"testing"

	. "github.com/longnguyen11288/nomadic-router/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "github.com/longnguyen11288/nomadic-router/Godeps/_workspace/src/github.com/onsi/gomega"
)

func TestGorb(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "gorb Test Suite")
}
