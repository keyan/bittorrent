package tracker

import (
	"errors"
	"net/http"
	"time"

	"github.com/google/go-querystring/query"

	"github.com/keyan/bittorrent/peers"
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
	Uploaded   string `url:"uploaded"`
	Downloaded string `url:"downloaded"`
	Left       string `url:"left"`
	Compact    int    `url:"compact"`
	NoPeerID   bool   `url:"no_peer_id"`
	Event      int    `url:"event"`
}

type Response struct {
	Peers    []peers.Peer
	Seeders  int
	Leechers int
}

func (t *Tracker) GetRequest(rp RequestParams) (*Response, error) {
	if t.nextAnnounceAfter > time.Now().Unix() {
		return nil, errors.New("cannot contact tracker yet")
	}

	if !t.hasBeenContacted {
		rp.event = TRACKER_STARTED_EVENT
	}

	v, _ := query.Values(rp)
	_, err := http.Get(t.url + "?" + v)
	if err != nil {
		return nil, err
	}

	t.hasBeenContacted = true

	r := Response{}

	return &r, nil
}

func New(url string) *Tracker {
	return &Tracker{
		url:               url,
		hasBeenContacted:  false,
		nextAnnounceAfter: 0,
	}
}
