package client

import (
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
}

// func (t *Torrent) RunTorrent() {
// 	trk := tracker.New(t.trackerUrl)
// 	p := tracker.RequestParams{
// 		Left:     t.piecesLeft,
// 		Compact:  0,
// 		NoPeerID: true,
// 	}

// 	resp, err := trk.GetRequest(p)
// 	if err != nil {
// 		fmt.Println("request failed")
// 	}
// 	fmt.Println(resp)
// }

// Download starts the Client and downloads the file provided as a Torrent.
// This function blocks until the entire file is downloaded or until an
// unrecoverable error is reached.
func (cl *Client) Download() error {
	return nil
}

// New creates a new Client and returns it.
func New(tor *torrent.Torrent, trk *tracker.Tracker) (*Client, error) {
	return &Client{
		torrent: tor,
		tracker: trk,
	}, nil
}
