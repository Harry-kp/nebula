package handshake

import (
	"errors"
	"io"
)

type HandShake struct {
	Pstr     string
	PeerID   [20]byte
	InfoHash [20]byte
}

func (h *HandShake) Serialize() []byte {
	// Check the format of the handshake https://blog.jse.li/posts/torrent/
	buffer := make([]byte, len(h.Pstr)+49)
	// Write length of Pstr on the index 0
	buffer[0] = byte(len(h.Pstr))
	offset := 1

	offset += copy(buffer[offset:], h.Pstr)
	offset += copy(buffer[offset:], make([]byte, 8))
	offset += copy(buffer[offset:], h.InfoHash[:])
	offset += copy(buffer[offset:], h.PeerID[:])
	return buffer
}

func Read(r io.Reader) (*HandShake, error) {
	// Read the length of the Pstr
	lengthBuf := make([]byte, 1)
	_, err := io.ReadFull(r, lengthBuf)
	if err != nil {
		return nil, err
	}
	// Read the Pstr Lenght
	pstrLength := int(lengthBuf[0])
	if pstrLength == 0 {
		return nil, errors.New("pstr length is 0")
	}

	// Read the Pstr
	hadnshakeBuffer := make([]byte, pstrLength+48)
	_, err = io.ReadFull(r, hadnshakeBuffer)
	if err != nil {
		return nil, err
	}

	// Parse the buffer
	h := &HandShake{}
	h.Pstr = string(hadnshakeBuffer[:pstrLength])
	copy(h.InfoHash[:], hadnshakeBuffer[pstrLength+8:pstrLength+28])
	copy(h.PeerID[:], hadnshakeBuffer[pstrLength+28:])
	return h, nil
}

func New(peerID, infoHash [20]byte) *HandShake {
	return &HandShake{
		Pstr:     "BitTorrent protocol",
		PeerID:   peerID,
		InfoHash: infoHash,
	}
}
