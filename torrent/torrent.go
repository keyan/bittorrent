package torrent

import (
	"crypto/sha1"

	"github.com/zeebo/bencode"
)

// Torrent is a utility type that abstracts the contents of the raw
// metainfo data to provide the most useful parts to other functions.
type Torrent struct {
	Name       string // The filename
	TrackerUrl string // The address for the tracker
	// SHA1 hash of the value of the 'info' key in metainfo. Not URL encoded
	// instead the Tracker will do that at request time.
	InfoHash       string
	PiecesToHash   map[int]string // Index i has the SHA1 hash for piece i
	BytesPerPiece  uint64
	PiecesAcquired uint64
	TotalPieces    uint64
}

// Piece represents a single BytesPerPiece sized unit of data downloaded from peers.
type Piece struct {
	Id   int // The piece number, index into PiecesToHash map
	Data []byte
}

// NewFromRawBytes converts the raw torrent data from file bytes to an
// internal Torrent type.
func NewFromRawBytes(fileBytes []byte) (*Torrent, error) {
	// Decode the torrent file to get the "Metainfo" map.
	var metainfo map[string]interface{}
	err := bencode.DecodeBytes(fileBytes, &metainfo)
	if err != nil {
		return nil, err
	}

	infoMap := metainfo["info"].(map[string]interface{})

	// InfoHash is used to set to the tracker, it must be the SHA1 hash of
	// the original bencoded 'info' dictionary value.
	bencodedInfo, err := bencode.EncodeBytes(infoMap)
	if err != nil {
		return nil, err
	}
	infoHashRaw := sha1.Sum(bencodedInfo)
	infoHash := string(infoHashRaw[:])

	// A string whose length is a multiple of 20. It is to be
	// subdivided into strings of length 20, each of which is the
	// SHA1 hash of the piece at the corresponding index
	piecesHash := infoMap["pieces"].(string)
	pieceMap := make(map[int]string)
	for i, j := 0, 0; i < len(piecesHash); i, j = i+20, j+1 {
		pieceMap[i] = piecesHash[i : i+20]
	}

	torrent := Torrent{
		Name:          infoMap["name"].(string),
		TrackerUrl:    metainfo["announce"].(string),
		InfoHash:      infoHash,
		PiecesToHash:  pieceMap,
		BytesPerPiece: uint64(infoMap["piece length"].(int64)),
		TotalPieces:   uint64(len(pieceMap)),
	}

	return &torrent, nil
}
