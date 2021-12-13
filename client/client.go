package client

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/keyan/bittorrent/peer"
	"github.com/keyan/bittorrent/torrent"
	"github.com/keyan/bittorrent/tracker"
)

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
	peers      []peer.Peer
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
			rp := tracker.RequestParams{
				InfoHash: cl.torrent.InfoHash,
				PeerID:   cl.peerID,
				Port:     cl.listenPort,
			}
			_, err := cl.tracker.GetRequest(rp)
			if err != nil {
				fmt.Println("Client: got error when pinging tracker")
			}

			// TODO: Reset so next tick respects what the Tracker instructed.
			tickDuration = 1 * time.Second
			ticker.Reset(tickDuration)
		}
	}
}

// Start begins the Client and downloads the file provided as a Torrent.
// This function blocks until the entire file is downloaded or until an
// unrecoverable error is reached.
func (cl *Client) Start() error {
	// TODO Open port to start listening
	// foo

	go cl.trackerLoop()

	// Enter main loop, continuously try to get pieces from any peers.
	for {
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
		torrent:    tor,
		tracker:    trk,
		peerID:     peerID,
		listenPort: 6881,
	}, nil
}
