package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/keyan/bittorrent/client"
	"github.com/keyan/bittorrent/torrent"
	"github.com/keyan/bittorrent/tracker"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	torrentFile := flag.String(
		"torrentFile", "example_data/debian.torrent", "The torrent to download")
	fileBytes, err := os.ReadFile(*torrentFile)
	check(err)

	// Abstract metainfo parsing, get a more useful struct that
	// has all the data we need.
	torrent, err := torrent.NewFromRawBytes(fileBytes)
	check(err)

	tracker, err := tracker.New(torrent.TrackerUrl)
	check(err)

	client, err := client.New(torrent, tracker)
	check(err)

	fmt.Printf(
		"Starting download for file: %s, total pieces: %d\n",
		torrent.Name, torrent.PiecesLeft)

	err = client.Start()
	check(err)

	fmt.Printf("Downloaded file %s\n", torrent.Name)
}
