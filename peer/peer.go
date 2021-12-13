package peer

import (
	"errors"
	"net"
)

type Peer struct {
	IP         net.IP
	Port       int16
	Connection *PeerConnection
}

type PeerConnection struct {
	amChoking      bool
	amInterested   bool
	peerChoking    bool
	peerInterested bool
}

func (p *Peer) Connect() error {
	if p.Connection != nil {
		return errors.New("already connected to peer")
	}

	// Choked and not-interested
	p.Connection = &PeerConnection{true, false, true, false}

	return nil
}
