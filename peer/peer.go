// package peer contains the structures used to manage the state of a Peer
// machine and is responsible for the logic around the Peer Wire Protocol
// and requesting missing Pieces from external peers.
package peer

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/keyan/bittorrent/torrent"
)

const connectTimeout = 2 * time.Second

// const (
// 	keepAlive msgType = iota
// 	choke
// 	unchoke
// 	interested
// 	notInterested
// 	have
// 	request
// 	piece
// 	cancel
// )

// var msgIDToType

type Peer struct {
	IP          net.IP
	Port        uint16
	State       *PeerState
	Conn        net.Conn
	PieceReqCh  chan int
	GivePieceCh chan *torrent.Piece
}

// PeerState indicates the current known state of an external Peer
// as well as what that Peer thinks of us (the client).
// We are choked when this Peer does not want to provide us with
// any Pieces, we are interested when we want to get Pieces from
// this Peer.
type PeerState struct {
	amChoking      bool
	amInterested   bool
	peerChoking    bool
	peerInterested bool
}

func (p *Peer) sendMessage(msgID int) error {
	// Messages in the protocol take the form of:
	//	<length prefix><message ID><payload>
	// The length prefix is a four byte big-endian value. The message
	// ID is a single decimal byte. The payload is message dependent.
	return nil
}

func (p *Peer) handshake(infoHash string, clientPeerID string) error {
	// The handshake is a required message and must be the first
	// message transmitted by the client. It is (49+len(pstr)) bytes long.
	// <pstrlen><pstr><reserved><info_hash><peer_id>
	buf := new(bytes.Buffer)
	// Length of the next string
	buf.WriteByte(19)
	// The default v1 protocol
	buf.WriteString("BitTorrent protocol")
	// 8 empty reserved bytes
	buf.Write(make([]byte, 8, 8))
	buf.WriteString(infoHash)
	buf.WriteString(clientPeerID)
	n, err := p.Conn.Write(buf.Bytes())
	if err != nil {
		return err
	}

	resp := make([]byte, 1024)
	n, err = p.Conn.Read(resp)
	if err != nil {
		return err
	}
	if n == 0 {
		return errors.New("handshake failed, got empty response")
	}

	// Start of infohash is 1 byte for the pstr len, the pstr, and 8 reserved bytes.
	infoHashOffset := uint8(resp[0]) + 1 + 8
	infoHashReceived := resp[infoHashOffset : infoHashOffset+20]
	if string(infoHashReceived) != infoHash {
		return errors.New("handshake failed, got incorrect infoHash")
	}

	return nil
}

// Connect opens a TCP connection with the external Peer, initiates the
// protocol required Handshake, then starts a loop to continuously ask
// for missing Pieces and send them back to the Client.
func (p *Peer) Connect(infoHash string, clientPeerID string) error {
	if p.State != nil {
		return errors.New("already connected to peer")
	}

	conn, err := net.DialTimeout(
		"tcp",
		fmt.Sprintf("%s:%v", p.IP.String(), p.Port),
		connectTimeout,
	)
	if err != nil {
		fmt.Printf("Peer: failed to connect to peer, %v\n", err)
		return err
	}

	p.Conn = conn
	// Choked and not-interested is the default start state.
	p.State = &PeerState{true, false, true, false}

	err = p.handshake(infoHash, clientPeerID)
	if err != nil {
		fmt.Printf(
			"Peer: peer %s failed to handshake, %v\n",
			p.IP.String(),
			err,
		)
		return err
	}
	fmt.Printf("Peer: handshake with peer %s was successful!\n", p.IP.String())

	return nil
}
