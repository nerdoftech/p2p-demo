package util

import (
	"crypto/ecdsa"
	"errors"
	"math/rand"
	"net"
	"time"

	"github.com/ethereum/go-ethereum/p2p/enode"

	"github.com/ethereum/go-ethereum/crypto"

	zlog "github.com/rs/zerolog/log"
)

// Tag logs with package name
var log = zlog.With().Str("pkg", "util").Logger()

// Is a wrapped LocalNode so we can get access to the key for later use, key is unexported in enode.LocalNode
type LNodeWrapper struct {
	Key       *ecdsa.PrivateKey
	DB        *enode.DB
	LocalNode *enode.LocalNode
}

// NewLNodeWrapper creates a wrapped LocalNode
func NewLNodeWrapper(path string) (*LNodeWrapper, error) {
	nw := &LNodeWrapper{}
	var err error
	nw.Key, err = GenerateNodeKey()
	if err != nil {
		return nil, err
	}
	err = nw.createDB(path)
	if err != nil {
		return nil, err
	}
	nw.LocalNode = enode.NewLocalNode(nw.DB, nw.Key)
	return nw, nil
}

// Creates a node db, file can be blank
func (l *LNodeWrapper) createDB(path string) error {
	log.Debug().Msg("opening DB")
	db, err := enode.OpenDB(path)
	if err != nil {
		msg := "failed to open DB"
		log.Error().Err(err).Msg(msg)
		return errors.New(msg)
	}
	l.DB = db
	return nil
}

func GenerateNodeKey() (*ecdsa.PrivateKey, error) {
	log.Debug().Msg("generating node key")
	key, err := crypto.GenerateKey()
	if err != nil {
		msg := "Failed to generate key"
		log.Error().Err(err).Msg(msg)
		return nil, errors.New(msg)
	}
	log.Debug().Msg("key generated")
	return key, nil
}

// Generate port number above 32767
func GenerateRandomPort() int {
	rand.Seed(time.Now().UnixNano())
	rn := rand.Int()
	numPorts := 65535 / 2

	return numPorts + (rn % numPorts)
}

// Create a new UDP listener
func NewUDPConn(ipStr string, port int) (*net.UDPConn, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		log.Error().Str("IP", ipStr).Msg("Could not parse IP string")
		return nil, errors.New("bad IP string")
	}
	addr := &net.UDPAddr{
		IP:   ip,
		Port: port,
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		msg := "can't start UDP listener"
		log.Error().Err(err).Msg(msg)
		return nil, errors.New(msg)
	}
	log.Debug().Interface("IP", addr.IP).Int("Port", addr.Port).Msg("starting UDP listener")
	return conn, nil
}
