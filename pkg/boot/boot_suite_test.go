package boot

import (
	"testing"

	"github.com/ethereum/go-ethereum/p2p/enode"

	"github.com/ethereum/go-ethereum/p2p/discover"

	"github.com/nerdoftech/p2p-demo/pkg/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
)

func TestBoot(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	RegisterFailHandler(Fail)
	RunSpecs(t, "Boot Suite")
}

var _ = Describe("boot node", func() {
	lhost := "127.0.0.1"
	Context("StartBootNode", func() {
		It("should work", func() {
			nw, err := util.NewLNodeWrapper("")
			Expect(err).Should(BeNil())

			cfg := &discover.Config{PrivateKey: nw.Key}
			port := util.GenerateRandomPort()
			disc, _, err := StartBootNode(nw.LocalNode, cfg, lhost, port)
			Expect(err).Should(BeNil())

			disc.Close()
		})
		It("fail for already allocated port", func() {
			port := util.GenerateRandomPort()
			conn, err := util.NewUDPConn(lhost, port)
			Expect(err).Should(BeNil())

			_, _, err = StartBootNode(nil, nil, lhost, port)
			Expect(err).Should(HaveOccurred())

			conn.Close()
		})
		It("fail for missing boot node IP", func() {
			port := util.GenerateRandomPort()
			nw, _ := util.NewLNodeWrapper("")
			cfg := &discover.Config{
				Bootnodes: []*enode.Node{
					&enode.Node{},
				},
			}
			_, _, err := StartBootNode(nw.LocalNode, cfg, lhost, port)
			Expect(err).Should(HaveOccurred())
		})
	})
})
