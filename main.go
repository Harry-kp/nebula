package main

import (
	"log"
	"os"

	"github.com/Harry-kp/nebula/torrentfile"
)

func main() {
	// Read args from terminal
	inPath := os.Args[1]
	outPath := os.Args[2]

	tf, err := torrentfile.Open(inPath)
	if err != nil {
		log.Fatal(err)
	}

	err = tf.DownloadToFile(outPath)
	if err != nil {
		log.Fatal(err)
	}
}
