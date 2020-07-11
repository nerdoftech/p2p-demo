package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/nerdoftech/p2p-demo/pkg/boot"
	"github.com/nerdoftech/p2p-demo/pkg/util"

	"github.com/ethereum/go-ethereum/p2p/discover"
	"github.com/ethereum/go-ethereum/p2p/netutil"

	log "github.com/sirupsen/logrus"
)

const (
	LOCAL_HOST = "127.0.0.1"
	LOCAL_NET  = "127.0.0.0/8"
)

func main() {
	log.SetLevel(log.DebugLevel)

	nw, err := util.NewLNodeWrapper("")
	if err != nil {
		log.WithError(err).Fatal("could not create local node")
	}

	netlist, err := netutil.ParseNetlist(LOCAL_NET)
	if err != nil {
		log.WithError(err).Fatal("could not parse netlist")
	}
	cfg := &discover.Config{
		PrivateKey:  nw.Key,
		NetRestrict: netlist,
	}

	port := util.GenerateRandomPort()
	disc, node, err := boot.StartBootNode(nw.LocalNode, cfg, LOCAL_HOST, port)
	if err != nil {
		log.WithError(err).Fatal()
	}

	// Send the node URL to stdout
	fmt.Println(node.URLv4())

	log.Debug("server started, waiting for signal to shutdown")
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	sig := <-ch
	log.WithField("signal", sig).Info("got signal, shutting down")

	disc.Close()
}
