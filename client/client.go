package client

import (
	"errors"
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/keyan/bittorrent/peer"
	"github.com/keyan/bittorrent/torrent"
	"github.com/keyan/bittorrent/tracker"
)

const MAX_PEERS int = 30

// Client is the top level object which manages a torrent download.
// It is responsible for contacting the tracker, keeping an active
// peer list and facilitating the peer protocol.
//
// For now a single Client manages a single Torrent, but in the
// future Client could be extended to handle multiple Torrents
// or the user of Client can create a Client for each Torrent
// to be downloaded in parallel.
type Client struct {
	torrent *torrent.Torrent
	tracker *tracker.Tracker
	// 20-byte string used as a unique ID for the client, generated
	// by the client at startup. Don't URL encode this, the Tracker
	// will do that.
	peerID     string
	listenPort int
	peers      map[string]peer.Peer
	// Used to share info about new peers to connect to.
	peersCh chan peer.Peer
	// Used by peers to send newly downloaded pieces from other peers.
	pieceCh chan *torrent.Piece
	// Used by peers to request which piece to ask for next. Messages
	// refer to the 0-indexed piece from the Torrent.
	wantPieceCh   chan int
	pieces        []*torrent.Piece
	missingPieces uint64
}

// trackerLoop regularly pings the tracker to inform it of our progress and
// collect an updated peers list.
func (cl *Client) trackerLoop() {
	// Tick immediately at first.
	tickDuration := 1 * time.Millisecond
	ticker := time.NewTicker(tickDuration)

	for {
		select {
		case <-ticker.C:
			ticker.Stop()

			rp := tracker.RequestParams{
				InfoHash: cl.torrent.InfoHash,
				PeerID:   cl.peerID,
				Port:     cl.listenPort,
			}
			resp, err := cl.tracker.GetRequest(rp)
			if err != nil {
				fmt.Println("Client: got error when pinging tracker")
			}

			fmt.Printf(
				"Client: received %d peers from tracker\n",
				len(resp.Peers))

			for _, p := range resp.Peers {
				cl.peersCh <- p
			}

			// Reset so next tick respects what the Tracker instructed
			tickDuration = time.Duration(resp.IntervalSecs) * time.Second
			ticker.Reset(tickDuration)
			fmt.Printf(
				"Client: next tracker request in %d secs\n",
				resp.IntervalSecs)
		}
	}
}

func (cl *Client) managePieces() {
	go func() {
		for missingPieces != 0 {
			// TODO Could be optimized, for now simplify how we assign
			// new pieces to request by just finding the next missing
			// one.
			for pieceIdx, piece := range cl.pieces {
				if piece == nil {
					cl.wantPieceCh <- pieceIdx
				}
			}
		}
		fmt.Println("Client: downloaded all pieces")
	}()

	for cl.missingPieces != 0 {
		select {
		case newPiece := <-cl.givePieceCh:
			// No need to lock access, we only ever write to one
			// index at one time. We might be reading this index
			// at the same time, but it doesn't matter.
			if cl.pieces[newPiece.Id] == nil {
				cl.pieces[newPiece.Id] = newPiece
				cl.missingPieces--
			}
		}
	}

}

// managePeers is responsible for:
//	1. receiving new peer information
//	2. starting peer connections
//
// TODO Should prune peers that stop responding
func (cl *Client) managePeers() {
	for {
		select {
		case newPeer := <-cl.peersCh:
			if len(cl.peers) >= MAX_PEERS {
				fmt.Println("Client: ignoring new peers, already at max")
				continue
			}
			// Already connected to this Peer.
			if _, ok := cl.peers[newPeer.IP.String()]; ok {
				continue
			}
			newPeer.PieceReqCh = cl.wantPieceCh
			newPeer.GivePieceCh = cl.givePieceCh
			err := newPeer.Connect(cl.torrent.InfoHash, cl.peerID)
			if err != nil {
				continue
			}
			cl.peers[newPeer.IP.String()] = newPeer
		}
	}
}

// Start begins the Client and downloads the file provided as a Torrent.
// This function blocks until the entire file is downloaded or until an
// unrecoverable error is reached.
func (cl *Client) Start() error {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", cl.listenPort))
	if err != nil {
		return err
	}

	// Routinely contact tracker and get updated peer list.
	go cl.trackerLoop()
	// Start all peers and manage new peer information passed via peersCh.
	go cl.managePeers()
	// Manage piece requests, answering Peer queries for which Piece to get
	// next.
	go cl.managePieces()

	// Continuously accept new peer messages and process them in serial.
	// For now not doing concurrent connections, let's see if this is
	// enough to not get choked by peers.
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf(
				"Client: got error when accepting inbound request, %w\n",
				err)
			continue
		}
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Printf(
				"Client: got error when reading inbound request, %w\n",
				err)
			continue
		}
		fmt.Printf("Read %d bytes\n", n)
		fmt.Println(buf[:n])
	}
	return nil
}

// New creates a new Client and returns it.
func New(tor *torrent.Torrent, trk *tracker.Tracker) (*Client, error) {
	// Some unique 20 char sequence to identify this Client.
	randBytes := make([]byte, 20)
	n, err := rand.Read(randBytes)
	if n != 20 || err != nil {
		return nil, errors.New("Couldn't generate peerID, this shouldn't happen")
	}
	peerID := string(randBytes)

	return &Client{
		torrent:       tor,
		tracker:       trk,
		peerID:        peerID,
		listenPort:    6881,
		peers:         make(map[string]peer.Peer),
		peersCh:       make(chan peer.Peer),
		wantPieceCh:   make(chan int),
		pieces:        make([]*torrent.Piece, int(tor.TotalPieces)),
		missingPieces: tor.TotalPieces,
		givePieceCh:   make(chan *torrent.Piece),
	}, nil
}
