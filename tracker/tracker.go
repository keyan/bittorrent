// package Tracker is used to manage a connection to a single Tracker host
// for a Torrent. A tracker is a single IP which is responsible for providing
// peer lists to clients and collecting download/upload information about
// torrents.
package tracker

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/go-querystring/query"

	"github.com/keyan/bittorrent/peer"
)

const (
	TRACKER_STARTED_EVENT   = "started"
	TRACKER_COMPLETED_EVENT = "completed"
	TRACKER_STOPPED_EVENT   = "stopped"
)

type Tracker struct {
	url               string
	hasBeenContacted  bool
	nextAnnounceAfter int64
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

type Response struct {
	Peers    []peer.Peer
	Seeders  int
	Leechers int
}

func (t *Tracker) GetRequest(rp RequestParams) (*Response, error) {
	if t.nextAnnounceAfter > time.Now().Unix() {
		return nil, errors.New("cannot contact tracker yet")
	}

	if !t.hasBeenContacted {
		rp.Event = TRACKER_STARTED_EVENT
		t.hasBeenContacted = true
	}

	v, _ := query.Values(rp)
	resp, err := http.Get(t.url + "?" + v.Encode())
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	fmt.Println(resp)

	if resp.StatusCode != 200 {
		return nil, errors.New("Received non-success response")
	}

	r := Response{}

	return &r, nil
}

func New(url string) (*Tracker, error) {
	return &Tracker{
		url:               url,
		hasBeenContacted:  false,
		nextAnnounceAfter: 0,
	}, nil
}
