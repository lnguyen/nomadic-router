package leafnodes_test

import (
	. "github.com/longnguyen11288/nomadic-router/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "github.com/longnguyen11288/nomadic-router/Godeps/_workspace/src/github.com/onsi/gomega"

	"testing"
)

func TestLeafNode(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "LeafNode Suite")
}
