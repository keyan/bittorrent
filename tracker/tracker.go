package tracker

import (
	"errors"
	"net/http"
	"time"

	"github.com/keyan/bittorrent/bencode"
	"github.com/keyan/bittorrent/types"
)

const (
	TRACKER_STARTED_EVENT   = "started"
	TRACKER_COMPLETED_EVENT = "completed"
	TRACKER_STOPPED_EVENT   = "stopped"
)

type Tracker struct {
	url               string
	hasBeenContacted  bool
	nextAnnounceAfter int
}

type Request struct {
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

type Response struct {
	Peers    []types.Peer
	Seeders  int
	Leechers int
}

func (t *Tracker) GetRequest(Request) (*Response, error) {
	if t.nextAnnounceAfter > time.Now().Unix() {
		return nil, errors.New("cannot contact tracker yet")
	}

	resp, err := http.Get(t.url)
	if err != nil {
		return nil, err
	}

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
