package boot

import (
	"errors"
	"net"

	"github.com/ethereum/go-ethereum/p2p/discover"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/nerdoftech/p2p-demo/pkg/util"
	zlog "github.com/rs/zerolog/log"
)

// Tag logs with package name
var log = zlog.With().Str("pkg", "bootnode").Logger()

func StartBootNode(ln *enode.LocalNode, cfg *discover.Config, ip string, port int) (*discover.UDPv4, *enode.Node, error) {
	errMsg := "failed to start server"

	log.Debug().Msg("starting UDP server")
	conn, err := util.NewUDPConn(ip, port)
	if err != nil {
		log.Err(err).Msg(errMsg)
		return nil, nil, errors.New(errMsg)
	}

	log.Debug().Msg("starting p2p discovery server")
	dUDP, err := discover.ListenV4(conn, ln, *cfg)
	if err != nil {
		log.Err(err).Msg(errMsg)
		return nil, nil, errors.New(errMsg)
	}

	// Needed to get node URL
	ipAddr := conn.LocalAddr().(*net.UDPAddr).IP
	node := enode.NewV4(&cfg.PrivateKey.PublicKey, ipAddr, 0, port)

	return dUDP, node, nil
}
