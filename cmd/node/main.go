package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ethereum/go-ethereum/p2p"

	"github.com/nerdoftech/p2p-demo/pkg/node"
	"github.com/nerdoftech/p2p-demo/pkg/util"

	log "github.com/sirupsen/logrus"
)

const (
	LOCAL_HOST = "127.0.0.1"
	LOCAL_NET  = "127.0.0.0/8"
	bootnode   = "enode://d794480cbb67ea51c1160c45ac6a569fa8d2f939584fe5c68fb078bbf93e1da0523aae463cd2fecb791b8d265a138d76c602b587d036614ba6073c1e5810304f@127.0.0.1:0?discport=42640"
)

func main() {
	log.SetLevel(log.DebugLevel)

	addr := fmt.Sprintf("%s:%d", LOCAL_HOST, util.GenerateRandomPort())

	s, err := node.NewP2PServer("joe", bootnode, LOCAL_NET, addr)
	if err != nil {
		log.WithError(err).Fatal("failed to create new server")
	}

	eventChan := make(chan *p2p.PeerEvent, 10)
	go printEvents(eventChan)
	s.SubscribeEvents(eventChan)

	log.Debug("starting server")
	err = s.Start()
	if err != nil {
		log.WithError(err).Fatal("Unable to start server")
	}

	log.Debug("server started, waiting for signal to shutdown")
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	sig := <-ch
	log.WithField("signal", sig).Info("got signal, shutting down")

	s.Stop()
}

func printEvents(ch chan *p2p.PeerEvent) {
	for {
		event := <-ch
		log.
			WithField("type", event.Type).
			WithField("peer", event.Peer).
			WithField("local_addr", event.LocalAddress).
			WithField("remote_addr", event.RemoteAddress).
			WithField("error", event.Error).
			Info("received new event")
	}
}
