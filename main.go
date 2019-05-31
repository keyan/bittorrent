package main

import (
	"fmt"
	"os"

	"github.com/keyan/bittorrent/bencode"
)

/*
all integers in the peer wire protocol are encoded as four byte big-endian values. This includes the length prefix on all messages that come after the handshake.
*/

type ClientConnection struct {
	amChoking      bool
	amInterested   bool
	peerChoking    bool
	peerInterested bool
}

const (
	TRACKER_STARTED_EVENT   = "started"
	TRACKER_COMPLETED_EVENT = "completed"
	TRACKER_STOPPED_EVENT   = "stopped"
)

type TrackerRequest struct {
	info_hash  string
	peer_id    string
	port       int //6881
	uploaded   string
	downloaded string
	left       string
	compact    int
	no_peer_id bool
	event      int
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func openConnection() (*ClientConnection, error) {
	// Choked and not-interested
	return ClientConnection{1, 0, 1, 0}
}

func main() {
	f, err := os.Open("example.torrent")
	check(err)

	metainfo, err := bencode.Decode(f)
	check(err)

	fmt.Println(metainfo)
}
