package peers

import (
	"encoding/binary"
	"fmt"
	"net"
)

// fetch peers from the tracker
type Peer struct {
	IP   net.IP
	Port uint16
}

func Unmarshal(peersBytes []byte) ([]Peer, error) {
	const peerSize = 6
	if len(peersBytes)%peerSize != 0 {
		return nil, fmt.Errorf("Invalid Peers Response")
	}
	totalPeers := len(peersBytes) / peerSize
	peers := make([]Peer, totalPeers)
	for i := 0; i < totalPeers; i++ {
		startIdx := i * peerSize
		peers[i].IP = net.IP(peersBytes[startIdx : startIdx+4])
		peers[i].Port = binary.BigEndian.Uint16(peersBytes[startIdx+4 : startIdx+6])
	}
	return peers, nil
}

func (p *Peer) String() string {
	return fmt.Sprintf("%s:%d", p.IP, p.Port)
}
