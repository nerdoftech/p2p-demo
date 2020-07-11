package node

import (
	"errors"

	elog "github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/p2p/netutil"

	"github.com/nerdoftech/p2p-demo/pkg/util"

	log "github.com/sirupsen/logrus"
)

func NewP2PServer(name string, bootnode string, rstNets string, addr string) (*p2p.Server, error) {

	nl, err := netutil.ParseNetlist(rstNets)
	if err != nil {
		msg := "failed to parse netlist"
		log.WithError(err).Error(msg)
		return nil, errors.New(msg)
	}
	log.WithField("allowed nets", nl).Debug("setting allowed nets")

	bn, err := enode.Parse(enode.ValidSchemes, bootnode)
	if err != nil {
		msg := "could not parse boot node"
		log.WithError(err).Error(msg)
		return nil, errors.New(msg)
	}

	key, err := util.GenerateNodeKey()
	if err != nil {
		log.WithError(err).Error()
		return nil, err
	}

	s := &p2p.Server{
		Config: p2p.Config{
			MaxPeers:       50,
			PrivateKey:     key,
			Name:           name,
			BootstrapNodes: []*enode.Node{bn},
			NetRestrict:    nl,
			ListenAddr:     addr,
			Logger:         elog.Root(),
			NodeDatabase:   "",
		},
	}

	return s, nil
}
