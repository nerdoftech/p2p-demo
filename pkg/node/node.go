package node

import (
	"errors"
	"os"

	"github.com/rs/zerolog"

	elog "github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/p2p/netutil"

	"github.com/nerdoftech/p2p-demo/pkg/util"

	zlog "github.com/rs/zerolog/log"
)

// Tag logs with package name
var log = zlog.With().Str("pkg", "node").Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr})

// NewP2PServer returns a configured p2p node
func NewP2PServer(name string, bootnode string, rstNets string, addr string) (*p2p.Server, error) {
	nl, err := netutil.ParseNetlist(rstNets)
	if err != nil {
		msg := "failed to parse netlist"
		log.Error().Err(err).Msg(msg)
		return nil, errors.New(msg)
	}
	log.Debug().Interface("allowed nets", nl).Msg("setting allowed nets")

	bn, err := enode.Parse(enode.ValidSchemes, bootnode)
	if err != nil {
		msg := "could not parse boot node"
		log.Error().Err(err).Msg(msg)
		return nil, errors.New(msg)
	}
	log.Debug().Interface("ID", bn.ID()).Str("url", bn.String()).Msg("Parsed boot node")

	key, err := util.GenerateNodeKey()
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	// P2P server logs to logrus
	p2pLog := elog.Root()
	handler := &p2pHandler{
		// tag p2p server logs
		logger: zlog.With().Str("pkg", "p2p-server").Str("node_name", name).Logger(),
	}
	p2pLog.SetHandler(handler)

	s := &p2p.Server{
		Config: p2p.Config{
			MaxPeers:       50,
			PrivateKey:     key,
			Name:           name,
			BootstrapNodes: []*enode.Node{bn},
			NetRestrict:    nl,
			ListenAddr:     addr,
			Logger:         p2pLog,
			NodeDatabase:   "",
		},
	}
	return s, nil
}

// Wrap zerolog as handler
type p2pHandler struct {
	logger zerolog.Logger
}

// Log makes p2pHandler match elog.Handler interface
func (p *p2pHandler) Log(rec *elog.Record) error {
	switch rec.Lvl {
	case elog.LvlTrace:
		p.logger.Trace().Interface("ctx", rec.Ctx).Msg(rec.Msg)
	case elog.LvlDebug:
		p.logger.Debug().Interface("ctx", rec.Ctx).Msg(rec.Msg)
	case elog.LvlInfo:
		p.logger.Info().Interface("ctx", rec.Ctx).Msg(rec.Msg)
	case elog.LvlWarn:
		p.logger.Warn().Interface("ctx", rec.Ctx).Msg(rec.Msg)
	case elog.LvlCrit:
		p.logger.Fatal().Interface("ctx", rec.Ctx).Msg(rec.Msg)
	}
	return nil
}
