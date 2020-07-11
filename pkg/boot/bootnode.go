package boot

import (
	"errors"
	"net"

	"github.com/ethereum/go-ethereum/p2p/discover"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/nerdoftech/p2p-demo/pkg/util"
	log "github.com/sirupsen/logrus"
)

func StartBootNode(ln *enode.LocalNode, cfg *discover.Config, ip string, port int) (*discover.UDPv4, *enode.Node, error) {
	errMsg := "failed to start server"

	log.Debug("starting UDP server")
	conn, err := util.NewUDPConn(ip, port)
	if err != nil {
		log.WithError(err).Error(errMsg)
		return nil, nil, errors.New(errMsg)
	}

	log.Debug("starting p2p discovery server")
	dUDP, err := discover.ListenV4(conn, ln, *cfg)
	if err != nil {
		log.WithError(err).Error(errMsg)
		return nil, nil, errors.New(errMsg)
	}

	ipAddr := conn.LocalAddr().(*net.UDPAddr).IP
	node := enode.NewV4(&cfg.PrivateKey.PublicKey, ipAddr, 0, port)
	return dUDP, node, nil
}
