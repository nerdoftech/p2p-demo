package util

import (
	"crypto/ecdsa"
	"errors"
	"math/rand"
	"net"
	"time"

	"github.com/ethereum/go-ethereum/p2p/enode"

	"github.com/ethereum/go-ethereum/crypto"

	log "github.com/sirupsen/logrus"
)

type LNodeWrapper struct {
	Key       *ecdsa.PrivateKey
	DB        *enode.DB
	LocalNode *enode.LocalNode
}

func NewLNodeWrapper(path string) (*LNodeWrapper, error) {
	ln := &LNodeWrapper{}
	err := ln.generateNodeKey()
	if err != nil {
		return nil, err
	}
	err = ln.createDB(path)
	if err != nil {
		return nil, err
	}
	ln.LocalNode = enode.NewLocalNode(ln.DB, ln.Key)
	return ln, nil
}

// Creates a node db, file can be blank
func (l *LNodeWrapper) createDB(path string) error {
	log.Debug("opening DB")
	db, err := enode.OpenDB(path)
	if err != nil {
		msg := "failed to open DB"
		log.WithError(err).Error(msg)
		return errors.New(msg)
	}
	l.DB = db
	return nil
}

func (l *LNodeWrapper) generateNodeKey() error {
	var err error
	l.Key, err = GenerateNodeKey()
	return err
}

func GenerateNodeKey() (*ecdsa.PrivateKey, error) {
	log.Debug("generating node key")
	key, err := crypto.GenerateKey()
	if err != nil {
		msg := "Failed to generate key"
		log.WithError(err).Error(msg)
		return nil, errors.New(msg)
	}
	log.Debug("key generated")
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
		log.WithField("IP", ipStr).Error("Could not parse IP string")
		return nil, errors.New("bad IP string")
	}
	addr := &net.UDPAddr{
		IP:   ip,
		Port: port,
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		msg := "can't start UDP listener"
		log.WithError(err).Error(msg)
		return nil, errors.New(msg)
	}
	log.WithField("IP", addr.IP).WithField("Port", addr.Port).Debug("starting UDP listener")
	return conn, nil
}
