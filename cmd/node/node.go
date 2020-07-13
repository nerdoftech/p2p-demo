package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ethereum/go-ethereum/p2p"

	"github.com/nerdoftech/p2p-demo/pkg/node"
	"github.com/nerdoftech/p2p-demo/pkg/util"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	LOCAL_HOST   = "127.0.0.1"
	LOCAL_NET    = "127.0.0.0/8"
	DEFAULT_PORT = 30303
)

var (
	flgBoot       = flag.String("bootnode", "", "The URL of the desired bootnode (required)")
	flgAddr       = flag.String("addr", LOCAL_HOST, "The local IP address for the bootnode discovery server")
	flgNetlist    = flag.String("netlist", LOCAL_NET, "A comma seperated list of allowed networks in CIDR format")
	flgPort       = flag.Int("port", DEFAULT_PORT, "UDP port for ")
	flgName       = flag.String("name", "", "Name of the node, default is node-port")
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

	if *flgRandomPort {
		*flgPort = util.GenerateRandomPort()
	}
	addr := fmt.Sprintf("%s:%d", *flgAddr, *flgPort)

	if *flgName == "" {
		*flgName = fmt.Sprintf("node-%d", *flgPort)
	}
	s, err := node.NewP2PServer(*flgName, *flgBoot, *flgNetlist, addr)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create new server")
	}

	eventChan := make(chan *p2p.PeerEvent, 10)
	go printEvents(eventChan)
	s.SubscribeEvents(eventChan)

	log.Debug().Msg("starting server")
	err = s.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to start server")
	}

	log.Debug().Msg("server started, waiting for signal to shutdown")
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	sig := <-ch
	log.Info().Interface("signal", sig).Msg("got signal, shutting down")

	s.Stop()
}

// Print server events
func printEvents(ch chan *p2p.PeerEvent) {
	for {
		event := <-ch
		logEvent := log.Info().
			Interface("type", event.Type).
			Interface("peer", event.Peer).
			Str("local_addr", event.LocalAddress).
			Str("remote_addr", event.RemoteAddress)

		// Only add error field if we have text
		if event.Error != "" {
			logEvent = logEvent.Str("error", event.Error)
		}

		logEvent.Msg("received new event")
	}
}
