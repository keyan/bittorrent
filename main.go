package main

import (
	"fmt"
	"os"

	"github.com/keyan/bittorrent/bencode"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	f, err := os.Open("example.torrent")
	check(err)

	metainfo, err := bencode.Decode(f)
	check(err)

	fmt.Println(metainfo)
}
