package torrent

import (
	"crypto/sha1"

	"github.com/keyan/bittorrent/bencode"
)

// Torrent is a utility type that abstracts the contents of the raw
// metainfo data to provide the most useful parts to other functions.
type Torrent struct {
	Name           string         // The filename
	TrackerUrl     string         // The address for the tracker
	InfoHash       string         // SHA1 hash of the value of the 'info' key in metainfo
	PiecesToHash   map[int]string // Index i has the SHA1 hash for piece i
	BytesPerPiece  uint64
	PiecesAcquired uint64
	PiecesLeft     uint64
}

// NewFromRawMetainfo converts the raw metainfo map to a more useful
// internal Torrent type.
func NewFromRawMetainfo(metainfo map[string]interface{}) (*Torrent, error) {
	infoMap := metainfo["info"].(map[string]interface{})

	// InfoHash is used to set to the tracker, it must be the SHA1 hash of
	// the original bencoded 'info' dictionary value.
	bencodedInfo, err := bencode.Encode(infoMap)
	if err != nil {
		return nil, err
	}
	infoHash := h.Sum(bencodedInfo)

	// A string whose length is a multiple of 20. It is to be
	// subdivided into strings of length 20, each of which is the
	// SHA1 hash of the piece at the corresponding index
	piecesHash := infoMap["pieces"].(string)
	pieceMap := make(map[int]string)
	for i, j := 0, 0; i < len(piecesHash); i, j = i+20, j+1 {
		pieceMap[i] = piecesHash[i : i+20]
	}

	torrent := Torrent{
		Name:           infoMap["name"].(string),
		TrackerUrl:     metainfo["announce"].(string),
		InfoHash:       infoHash,
		PiecesToHash:   pieceMap,
		BytesPerPiece:  infoMap["piece length"].(uint64),
		PiecesAcquired: 0,
		PiecesLeft:     uint64(len(pieceMap)),
	}

	return &torrent, nil
}
