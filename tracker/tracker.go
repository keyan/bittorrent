// package Tracker is used to manage a connection to a single Tracker host
// for a Torrent. A tracker is a single IP which is responsible for providing
// peer lists to clients and collecting download/upload information about
// torrents.
package tracker

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/google/go-querystring/query"
	"github.com/zeebo/bencode"

	"github.com/keyan/bittorrent/peer"
)

const (
	TRACKER_STARTED_EVENT   = "started"
	TRACKER_COMPLETED_EVENT = "completed"
	TRACKER_STOPPED_EVENT   = "stopped"
)

type Tracker struct {
	url              string
	hasBeenContacted bool
}

type RequestParams struct {
	InfoHash   string `url:"info_hash"`
	PeerID     string `url:"peer_id"`
	Port       int    `url:"port"`
	Uploaded   uint64 `url:"uploaded"`
	Downloaded uint64 `url:"downloaded"`
	Left       uint64 `url:"left"`
	Compact    int    `url:"compact"`
	NoPeerID   bool   `url:"no_peer_id"`
	Event      string `url:"event"`
}

// Response is the exported type used by callers to access tracker information.
type Response struct {
	Peers        []peer.Peer
	IntervalSecs int64
}

// trackerResponse is used for deserializing the raw text tracker response
// from bencode.
type TrackerResponse struct {
	FailureReason string `bencode:"failure reason"`
	// Raw bytes for the peer data, needs to be decoded
	Peers        string `bencode:"peers"`
	IntervalSecs int64  `bencode:"interval"`
}

// peersFromRawBytes converts the raw tracker `peers` data to a slice of
// internal Peer structs. According to the spec this data is:
//
//	a string consisting of multiples of 6 bytes.
//	First 4 bytes are the IP address and last 2 bytes
//	are the port number. All in network (big endian) notation.
func peersFromRawBytes(rawPeers []byte) []peer.Peer {
	peerList := make([]peer.Peer, 0)
	for i := 0; i < len(rawPeers); i = i + 6 {
		ip := net.IPv4(rawPeers[i], rawPeers[i+1], rawPeers[i+2], rawPeers[i+3])

		var port uint16
		buf := bytes.NewReader(rawPeers[i+4 : i+6])
		err := binary.Read(buf, binary.BigEndian, &port)
		if err != nil {
			fmt.Printf("Tracker: got error during peer decoding, %w", err)
			continue
		}

		peerList = append(peerList, peer.Peer{IP: ip, Port: port})
	}

	return peerList
}

// GetRequest pings the tracker specified in the Torrent and collects
// a list of active Peers that should be used to download the remaining
// pieces.
func (t *Tracker) GetRequest(rp RequestParams) (*Response, error) {
	if !t.hasBeenContacted {
		rp.Event = TRACKER_STARTED_EVENT
		t.hasBeenContacted = true
	}

	v, _ := query.Values(rp)
	resp, err := http.Get(t.url + "?" + v.Encode())
	if err != nil {
		return nil, fmt.Errorf("Tracker: request failed, %w", err)

	}
	if resp.StatusCode != 200 {
		return nil, errors.New("Received non-success response")
	}

	// Decode response body
	defer resp.Body.Close()
	respData, err := io.ReadAll(resp.Body)
	var trkResp TrackerResponse
	bencode.DecodeBytes(respData, &trkResp)

	// This field is only set if the others are not.
	if len(trkResp.FailureReason) > 0 {
		return nil, fmt.Errorf(
			"Tracker: got failure response, %s", trkResp.FailureReason)
	}

	r := Response{
		Peers:        peersFromRawBytes([]byte(trkResp.Peers)),
		IntervalSecs: trkResp.IntervalSecs,
	}

	return &r, nil
}

func New(url string) (*Tracker, error) {
	return &Tracker{
		url:              url,
		hasBeenContacted: false,
	}, nil
}
