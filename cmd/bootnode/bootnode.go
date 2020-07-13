package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/nerdoftech/p2p-demo/pkg/boot"
	"github.com/nerdoftech/p2p-demo/pkg/util"

	"github.com/ethereum/go-ethereum/p2p/discover"
	"github.com/ethereum/go-ethereum/p2p/netutil"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	LOCAL_HOST   = "127.0.0.1"
	LOCAL_NET    = "127.0.0.0/8"
	DEFAULT_PORT = 30303
)

var (
	flgAddr       = flag.String("addr", LOCAL_HOST, "The local IP address for the bootnode discovery server")
	flgNetlist    = flag.String("netlist", LOCAL_NET, "A comma seperated list of allowed networks in CIDR format")
	flgPort       = flag.Int("port", DEFAULT_PORT, "UDP port for ")
	flgRandomPort = flag.Bool("random", false, "Generates a random port for the server to listen on, overrides 'port' setting")
	flgLogLvl     = flag.String("log", "info", "Sets the log level")
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	flag.Parse()

	lvl, err := zerolog.ParseLevel(*flgLogLvl)
	if err != nil {
		log.Fatal().Err(err).Msg("could not parse log level")
	}
	zerolog.SetGlobalLevel(lvl)

	nw, err := util.NewLNodeWrapper("")
	if err != nil {
		log.Fatal().Err(err).Msg("could not create local node")
	}

	list, err := netutil.ParseNetlist(*flgNetlist)
	if err != nil {
		log.Fatal().Err(err).Msg("could not parse netlist")
	}
	cfg := &discover.Config{
		PrivateKey:  nw.Key,
		NetRestrict: list,
	}

	if *flgRandomPort {
		*flgPort = util.GenerateRandomPort()
	}
	disc, node, err := boot.StartBootNode(nw.LocalNode, cfg, *flgAddr, *flgPort)
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}

	// Send the node URL to stdout
	fmt.Println(node.URLv4())

	log.Debug().Msgf("server started, waiting for signal to shutdown")
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	sig := <-ch
	log.Info().Interface("signal", sig).Msg("got signal, shutting down")

	disc.Close()
}
