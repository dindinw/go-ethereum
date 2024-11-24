package p2p

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"net"
)

type QngServer struct {
	cfg *Config

	ourHandshake *protoHandshake
	peerFeed     event.Feed
	log          log.Logger
}

func (srv *QngServer) Start(s *Server) error {
	srv.cfg = &s.Config
	srv.ourHandshake = s.ourHandshake
	srv.log = s.log
	return nil
}

func (srv *QngServer) Connect(fd net.Conn, dialDest *enode.Node) (bool, error) {
	flags := trustedConn
	if dialDest == nil {
		flags |= inboundConn
	} else {
		flags |= dynDialedConn
	}
	return srv.SetupConn(fd, flags, dialDest)
}

func (srv *QngServer) SetupConn(fd net.Conn, flags connFlag, dialDest *enode.Node) (bool, error) {
	c := &conn{fd: fd, flags: flags, cont: make(chan error)}
	if dialDest == nil {
		c.transport = newRLPX(fd, nil)
	} else {
		c.transport = newRLPX(fd, dialDest.Pubkey())
	}

	ret, err := srv.setupConn(c, dialDest)
	if err != nil {
		if !c.is(inboundConn) {
			markDialError(err)
		}
		c.close(err)
	}
	return ret, err
}

func (srv *QngServer) setupConn(c *conn, dialDest *enode.Node) (bool, error) {
	// If dialing, figure out the remote public key.
	if dialDest != nil {
		dialPubkey := new(ecdsa.PublicKey)
		if err := dialDest.Load((*enode.Secp256k1)(dialPubkey)); err != nil {
			err = fmt.Errorf("%w: dial destination doesn't have a secp256k1 public key", errEncHandshakeError)
			srv.log.Trace("Setting up connection failed", "addr", c.fd.RemoteAddr(), "conn", c.flags, "err", err)
			return false, err
		}
	}

	// Run the RLPx handshake.
	remotePubkey, err := c.doEncHandshake(srv.cfg.PrivateKey)
	if err != nil {
		srv.log.Trace("Failed RLPx handshake", "addr", c.fd.RemoteAddr(), "conn", c.flags, "err", err)
		return false, fmt.Errorf("%w: %v", errEncHandshakeError, err)
	}
	if dialDest != nil {
		c.node = dialDest
	} else {
		c.node = nodeFromConn(remotePubkey, c.fd)
	}
	clog := srv.log.New("id", c.node.ID(), "addr", c.fd.RemoteAddr(), "conn", c.flags)

	// Run the capability negotiation handshake.
	phs, err := c.doProtoHandshake(srv.ourHandshake)
	if err != nil {
		clog.Trace("Failed p2p handshake", "err", err)
		return false, fmt.Errorf("%w: %v", errProtoHandshakeError, err)
	}
	if id := c.node.ID(); !bytes.Equal(crypto.Keccak256(phs.ID), id[:]) {
		clog.Trace("Wrong devp2p handshake identity", "phsid", hex.EncodeToString(phs.ID))
		return false, DiscUnexpectedIdentity
	}
	c.caps, c.name = phs.Caps, phs.Name
	err = srv.addPeerChecks(c)
	if err != nil {
		clog.Trace("Rejected peer", "err", err)
		return false, err
	}
	_, ret := srv.launchPeer(c)
	return ret, nil
}

func (srv *QngServer) addPeerChecks(c *conn) error {
	// Drop connections with no matching protocols.
	if len(srv.cfg.Protocols) > 0 && countMatchingProtocols(srv.cfg.Protocols, c.caps) == 0 {
		return DiscUselessPeer
	}
	return nil
}

func (srv *QngServer) launchPeer(c *conn) (*Peer, bool) {
	p := newPeer(srv.log, c, srv.cfg.Protocols)
	if srv.cfg.EnableMsgEvents {
		// If message events are enabled, pass the peerFeed
		// to the peer.
		p.events = &srv.peerFeed
	}
	ret := srv.runPeer(p)
	return p, ret
}

func (srv *QngServer) runPeer(p *Peer) bool {
	srv.peerFeed.Send(&PeerEvent{
		Type:          PeerEventTypeAdd,
		Peer:          p.ID(),
		RemoteAddress: p.RemoteAddr().String(),
		LocalAddress:  p.LocalAddr().String(),
	})

	// Run the per-peer main loop.
	remoteRequested, err := p.run()

	srv.peerFeed.Send(&PeerEvent{
		Type:          PeerEventTypeDrop,
		Peer:          p.ID(),
		Error:         err.Error(),
		RemoteAddress: p.RemoteAddr().String(),
		LocalAddress:  p.LocalAddr().String(),
	})

	return remoteRequested
}

func NewQngServer() *QngServer {
	qs := &QngServer{}
	return qs
}
