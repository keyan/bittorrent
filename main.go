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
	trk := tracker.New(t.trackerUrl)
	p := tracker.RequestParams{
		Left:     t.piecesLeft,
		Compact:  0,
		NoPeerID: true,
	}

	resp, err := trk.GetRequest(p)
	if err != nil {
		fmt.Println("request failed")
	}
	fmt.Println(resp)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	f, err := os.Open("torrents/flagfromserver.torrent")
	check(err)

	metainfo, err := bencode.Decode(f)
	check(err)

	infoMap := metainfo["info"].(map[string]interface{})

	torrent := Torrent{
		// name:          metainfo["title"].(string),
		trackerUrl:    metainfo["announce"].(string),
		piecesHash:    infoMap["pieces"].(string),
		bytesPerPiece: infoMap["piece length"].(uint64),
	}
	torrent.RunTorrent()
}
