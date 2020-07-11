package node

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestNode(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Node Suite")
}

const (
	bn = "enode://d794480cbb67ea51c1160c45ac6a569fa8d2f939584fe5c68fb078bbf93e1da0523aae463cd2fecb791b8d265a138d76c602b587d036614ba6073c1e5810304f@127.0.0.1:0?discport=42640"
)

var _ = Describe("node", func() {
	Context("NewP2PServer", func() {
		It("should work", func() {
			_, err := NewP2PServer("", bn, "", "127.0.0.1:55555")
			Expect(err).Should(BeNil())
		})
		It("fail because bad netlist", func() {
			_, err := NewP2PServer("", bn, "1234", "")
			Expect(err).Should(HaveOccurred())
		})
		It("fail because bad bootnode", func() {
			_, err := NewP2PServer("", "1234", "", "")
			Expect(err).Should(HaveOccurred())
		})
	})
})
