package main

import (
	"fmt"
	"os"

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
	fmt.Println(tf)
}
