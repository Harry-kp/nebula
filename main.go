package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/Harry-kp/bit-bandit/handshake"
	"github.com/Harry-kp/bit-bandit/torrentfile"
)

func main() {
	// Read args from terminal
	args := os.Args[1:]
	if len(args) != 1 {
		panic("Provide a torrent file path as the first argument")
	}
	// Open the file
	torrent_file_path := args[0]
	tf, err := torrentfile.Open(torrent_file_path)
	if err != nil {
		panic(err)
	}
	peerID := [20]byte{}
	port := uint16(1121)
	peerList, err := tf.FetchPeers(peerID, port)

	// Send and Recieve Handshake to peer
	conn, err := net.DialTimeout("tcp", peerList[0].String(), 3*time.Second)
	// Send the handshake to connection
	conn.Write(handshake.New(peerID, tf.InfoHash).Serialize())
	res, err := handshake.Read(conn)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
}
