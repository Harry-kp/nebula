package client

import (
	"bytes"
	"fmt"
	"net"
	"time"

	"github.com/Harry-kp/nebula/bitfield"
	"github.com/Harry-kp/nebula/handshake"
	"github.com/Harry-kp/nebula/message"
	"github.com/Harry-kp/nebula/peers"
)

type Client struct {
	Bitfield bitfield.Bitfield
	Conn     net.Conn
	Choked   bool
	infoHash [20]byte
	peerID   [20]byte
	peer     peers.Peer
}

func completeHandshake(conn net.Conn, infoHash, peerID [20]byte) (*handshake.HandShake, error) {
	conn.SetDeadline(time.Now().Add(3 * time.Second))
	defer conn.SetDeadline(time.Time{}) // Disable the

	sendHsk := handshake.New(peerID, infoHash)

	_, err := conn.Write(sendHsk.Serialize())
	if err != nil {
		return nil, err
	}

	resHsk, err := handshake.Read(conn)
	if err != nil {
		return nil, err
	}

	if !bytes.Equal(resHsk.InfoHash[:], infoHash[:]) {
		return nil, fmt.Errorf("Expected infohash %x, but got %x", infoHash, resHsk.InfoHash)
	}
	return resHsk, nil
}

func fetchBitField(conn net.Conn) (bitfield.Bitfield, error) {
	conn.SetDeadline(time.Now().Add(3 * time.Second))
	defer conn.SetDeadline(time.Time{})

	msg, err := message.Read(conn)
	if err != nil {
		return nil, err
	}

	if msg == nil {
		return nil, fmt.Errorf("Expected bitfield, but got nil")
	}

	if msg.ID != message.MsgBitfield {
		return nil, fmt.Errorf("Expected bitfield, but got ID %d", msg.ID)
	}

	return msg.Payload, nil
}

func New(peer peers.Peer, infoHash, peerID [20]byte) (*Client, error) {
	conn, err := net.DialTimeout("tcp", peer.String(), 3*time.Second)
	if err != nil {
		return nil, err
	}

	_, err = completeHandshake(conn, infoHash, peerID)
	if err != nil {
		conn.Close()
		return nil, err
	}

	bitfield, err := fetchBitField(conn)
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &Client{
		Bitfield: bitfield,
		Conn:     conn,
		Choked:   true,
		infoHash: infoHash,
		peerID:   peerID,
		peer:     peer,
	}, nil
}

func (c *Client) Read() (*message.Message, error) {
	msg, err := message.Read(c.Conn)
	if err != nil {
		return nil, err
	}
	return msg, nil

}

func (c *Client) SendRequest(index, begin, length int) error {
	msg := message.FormatRequest(index, begin, length)
	_, err := c.Conn.Write(msg.Serialize())
	return err
}

func (c *Client) SendInterested() error {
	msg := message.Message{ID: message.MsgInterested}
	_, err := c.Conn.Write(msg.Serialize())
	return err
}

func (c *Client) SendNotInterested() error {
	msg := message.Message{ID: message.MsgNotInterested}
	_, err := c.Conn.Write(msg.Serialize())
	return err
}

// SendUnchoke sends an Unchoke message to the peer
func (c *Client) SendUnchoke() error {
	msg := message.Message{ID: message.MsgUnchoke}
	_, err := c.Conn.Write(msg.Serialize())
	return err
}

// SendHave sends a Have message to the peer
func (c *Client) SendHave(index int) error {
	msg := message.FormatHave(index)
	_, err := c.Conn.Write(msg.Serialize())
	return err
}
