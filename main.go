package main

import (
	"fmt"
	"os"

	"github.com/keyan/bittorrent/bencode"
	"github.com/keyan/bittorrent/tracker"
)

type Torrent struct {
	name           string
	trackerUrl     string
	piecesHash     string
	bytesPerPiece  uint64
	data           []byte
	piecesAcquired uint64
	piecesLeft     uint64
}

func (t *Torrent) RunTorrent() {
	trk := tracker.New(torrent.trackerUrl)
	p := tracker.RequestParams{
		InfoHash:   nil,
		PeerID:     nil,
		Port:       nil,
		Uploaded:   nil,
		Downloaded: nil,
		Left:       t.piecesLeft,
		Compact:    0,
		NoPeerID:   1,
	}

	resp, err := trk.GetRequest(p)
	fmt.Println(resp)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	f, err := os.Open("torrents/example.torrent")
	check(err)

	metainfo, err := bencode.Decode(f)
	check(err)

	infoMap := metainfo["info"].(map[string]interface{})

	torrent := Torrent{
		name:          metainfo["title"].(string),
		trackerUrl:    metainfo["announce"].(string),
		piecesHash:    infoMap["pieces"].(string),
		bytesPerPiece: infoMap["piece length"].(uint64),
	}
	torrent.RunTorrent()
}
