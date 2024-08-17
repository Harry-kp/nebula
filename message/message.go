package message

import (
	"encoding/binary"
	"fmt"
	"io"
)

// The message id defines the type of message and it is stored in a single byte
type messageID uint8

const (
	MsgChoke         messageID = 0
	MsgUnchoke       messageID = 1
	MsgInterested    messageID = 2
	MsgNotInterested messageID = 3
	MsgHave          messageID = 4
	MsgBitfield      messageID = 5
	MsgRequest       messageID = 6
	MsgPiece         messageID = 7
	MsgCancel        messageID = 8
)

type Message struct {
	ID      messageID
	Payload []byte
}

// First We Ask the peer to get if they have the piece
func FormatHave(index int) *Message {
	m := Message{
		ID: MsgHave,
	}
	buffer := make([]byte, 4)
	binary.BigEndian.PutUint32(buffer, uint32(index))
	m.Payload = buffer
	return &m
}

// Then we parse the message to see if the peer has the piece
func ParseHave(msg *Message) (int, error) {
	if msg.ID != MsgHave {
		return 0, fmt.Errorf("Expected HAVE (ID %d), got ID %d", MsgHave, msg.ID)
	}

	if len(msg.Payload) != 4 {
		return 0, fmt.Errorf("Expected payload length 4, got %d", len(msg.Payload))
	}
	index := int(binary.BigEndian.Uint32(msg.Payload))
	return index, nil
}

func ParsePiece(index int, buf []byte, msg *Message) (int, error) {
	if msg.ID != MsgPiece {
		return 0, fmt.Errorf("Expected PIECE (ID %d), got ID %d", MsgPiece, msg.ID)
	}

	if len(msg.Payload) < 8 {
		return 0, fmt.Errorf("Expected payload length at least 8, got %d", len(msg.Payload))
	}

	pieceIndex := int(binary.BigEndian.Uint32(msg.Payload[0:4]))
	begin := int(binary.BigEndian.Uint32(msg.Payload[4:8]))

	if pieceIndex != index {
		return 0, fmt.Errorf("Expected piece index %d, got %d", index, pieceIndex)
	}
	data := msg.Payload[8:]
	if begin+len(data) > len(buf) {
		return 0, fmt.Errorf("Data too long [%d] for offset %d with length %d", len(data), begin, len(buf))
	}
	copy(buf[begin:], msg.Payload[8:])
	return len(data), nil
}

// FormatRequest creates a REQUEST message for the given piece if the peer has it
func FormatRequest(index, begin, length int) *Message {
	payload := make([]byte, 12)
	binary.BigEndian.PutUint32(payload[0:4], uint32(index))
	binary.BigEndian.PutUint32(payload[4:8], uint32(begin))
	binary.BigEndian.PutUint32(payload[8:12], uint32(length))
	return &Message{ID: MsgRequest, Payload: payload}
}

// Create a new message from the stream
func Read(r io.Reader) (*Message, error) {
	// Read the length of the message
	lengthBuf := make([]byte, 4)
	_, err := io.ReadFull(r, lengthBuf)
	if err != nil {
		return nil, err
	}
	length := binary.BigEndian.Uint32(lengthBuf)

	// keet the connection alive
	if length == 0 {
		return nil, nil
	}

	messageBuffer := make([]byte, length)
	_, err = io.ReadFull(r, messageBuffer)

	if err != nil {
		return nil, err
	}

	return &Message{
		ID:      messageID(messageBuffer[0]),
		Payload: messageBuffer[1:],
	}, nil
}

// Write the message to the stream
func (m *Message) Serialize() []byte {
	if m == nil {
		return make([]byte, 4)
	}
	length := uint32(len(m.Payload) + 1)
	buf := make([]byte, length+4)

	// Write the length of the message
	binary.BigEndian.PutUint32(buf[0:4], length)
	buf[4] = byte(m.ID)
	copy(buf[5:], m.Payload)
	return buf
}

func (m *Message) name() string {
	switch m.ID {
	case MsgChoke:
		return "choke"
	case MsgUnchoke:
		return "unchoke"
	case MsgInterested:
		return "interested"
	case MsgNotInterested:
		return "not interested"
	case MsgHave:
		return "have"
	case MsgBitfield:
		return "bitfield"
	case MsgRequest:
		return "request"
	case MsgPiece:
		return "piece"
	case MsgCancel:
		return "cancel"
	default:
		return fmt.Sprintf("Unknown#%d", m.ID)
	}
}

func (m *Message) String() string {
	return fmt.Sprintf("{ID: %s, Payload: [%v]}", m.name(), len(m.Payload))
}
