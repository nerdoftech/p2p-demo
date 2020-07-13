package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/p2p/enode"

	"github.com/nerdoftech/p2p-demo/pkg/boot"
	"github.com/nerdoftech/p2p-demo/pkg/node"
	"github.com/nerdoftech/p2p-demo/pkg/util"

	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/discover"
	"github.com/ethereum/go-ethereum/p2p/netutil"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	LOCAL_HOST = "127.0.0.1"
	LOCAL_NET  = "127.0.0.0/8"
)

func TestTest(t *testing.T) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	RegisterFailHandler(Fail)
	RunSpecs(t, "P2P Suite")
}

var _ = Describe("test p2p node interaction", func() {
	var discSrv *discover.UDPv4
	var bootnode string
	BeforeSuite(func() {
		// Start boot node
		nw, err := util.NewLNodeWrapper("")
		Expect(err).Should(BeNil())

		netlist, err := netutil.ParseNetlist(LOCAL_NET)
		Expect(err).Should(BeNil())
		cfg := &discover.Config{
			PrivateKey:  nw.Key,
			NetRestrict: netlist,
		}

		port := util.GenerateRandomPort()
		var node *enode.Node
		discSrv, node, err = boot.StartBootNode(nw.LocalNode, cfg, LOCAL_HOST, port)
		Expect(err).Should(BeNil())

		bootnode = node.URLv4()
		log.Info().Str("pkg", "test").Str("url", bootnode).Msg("allowing bootnode to start")
		// For some reason bootnode takes a while to startup properly
		time.Sleep(10 * time.Second)
	})
	AfterSuite(func() {
		discSrv.Close()
	})
	Context("run nodes", func() {
		// start 2 nodes, wait for peer add event and assert it matches the other node's ID.
		It("should work", func() {
			var peer1, peer2 enode.ID
			// Get node1 ready
			addr1 := fmt.Sprintf("%s:%d", LOCAL_HOST, util.GenerateRandomPort())
			node1, err := node.NewP2PServer("node1", bootnode, LOCAL_NET, addr1)
			Expect(err).Should(BeNil())

			// Get subscriber chan which will be used for assertion
			eventChan1 := make(chan *p2p.PeerEvent, 10)
			node1.SubscribeEvents(eventChan1)
			node1.Start()
			go waitForPeer(eventChan1, &peer1)

			time.Sleep(5 * time.Second)
			// Start node2
			addr2 := fmt.Sprintf("%s:%d", LOCAL_HOST, util.GenerateRandomPort())
			node2, err := node.NewP2PServer("node2", bootnode, LOCAL_NET, addr2)
			Expect(err).Should(BeNil())

			eventChan2 := make(chan *p2p.PeerEvent, 10)
			node2.SubscribeEvents(eventChan2)
			node2.Start()
			go waitForPeer(eventChan2, &peer2)

			// Wait for nodes to see each other as peers
			TIMEOUT := 45 * time.Second
			Eventually(
				func() enode.ID {
					return peer1
				},
				TIMEOUT, time.Second).
				Should(Equal(node2.LocalNode().ID()), "node1")
			Eventually(
				func() enode.ID {
					return peer2
				},
				TIMEOUT, time.Second).
				Should(Equal(node1.LocalNode().ID()), "node2")

			node1.Stop()
			node2.Stop()
		})
	})
})

// waits for peer add event and saves enode id to id pointer
func waitForPeer(ch chan *p2p.PeerEvent, id *enode.ID) {
	for {
		event := <-ch
		if event.Type == p2p.PeerEventTypeAdd {
			*id = event.Peer
		}
	}
}
