package torrentfile

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"os"

	"github.com/jackpal/bencode-go"
)

type bencodeInfo struct {
	Pieces      string `bencode:"pieces"`
	PieceLength int    `bencode:"piece length"`
	Length      int    `bencode:"length"`
	Name        string `bencode:"name"`
}

type bencodeTorrent struct {
	Announce string      `bencode:"announce"`
	Info     bencodeInfo `bencode:"info"`
}

type TorrentFile struct {
	Announce    string
	InfoHash    [20]byte
	PieceHashes [][20]byte
	PieceLength int
	Length      int
	Name        string
}

func (i *bencodeInfo) Hash() ([20]byte, error) {
	var buf bytes.Buffer
	err := bencode.Marshal(&buf, *i)
	if err != nil {
		return [20]byte{}, err
	}
	return sha1.Sum(buf.Bytes()), nil
}

func (info *bencodeInfo) splitPieceHashes() ([][20]byte, error) {
	hashLen := 20
	if len(info.Pieces)%hashLen != 0 {
		return nil, fmt.Errorf("invalid pieces length")
	}
	numHashes := len(info.Pieces) / hashLen
	hashes := make([][20]byte, numHashes)
	for i := 0; i < numHashes; i++ {
		copy(hashes[i][:], info.Pieces[i*hashLen:(i+1)*hashLen])
	}
	return hashes, nil
}

func Open(path string) (TorrentFile, error) {
	f, err := os.Open(path)
	if err != nil {
		return TorrentFile{}, err
	}
	defer f.Close()
	bto := bencodeTorrent{}
	err = bencode.Unmarshal(f, &bto)
	if err != nil {
		return TorrentFile{}, err
	}
	return bto.toTorrentFile()
}

func (bto *bencodeTorrent) toTorrentFile() (TorrentFile, error) {
	tf := TorrentFile{}
	hash, err := bto.Info.Hash()
	if err != nil {
		return tf, err
	}
	tf.InfoHash = hash
	tf.Announce = bto.Announce
	tf.PieceLength = bto.Info.PieceLength
	tf.Length = bto.Info.Length
	tf.Name = bto.Info.Name
	if piecesHash, err := bto.Info.splitPieceHashes(); err != nil {
		return tf, err
	} else {
		tf.PieceHashes = piecesHash
	}
	return tf, nil
}
