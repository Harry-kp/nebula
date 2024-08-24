package p2p

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"runtime"
	"time"

	"github.com/Harry-kp/nebula/client"
	"github.com/Harry-kp/nebula/logger"
	"github.com/Harry-kp/nebula/message"
	"github.com/Harry-kp/nebula/peers"
	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
)

// MaxBlockSize is the largest number of bytes a request can ask for
const maxBlockSize = 16384

// MaxBacklog is the number of unfulfilled requests a client can have in its pipeline
const maxBacklog = 5

type Torrent struct {
	Peers       []peers.Peer
	PeerID      [20]byte
	InfoHash    [20]byte
	PieceHashes [][20]byte
	PieceLength int
	Length      int
	Name        string
}

type pieceWork struct {
	index  int
	hash   [20]byte
	length int
}

type pieceResult struct {
	index int
	buf   []byte
}

type pieceProgress struct {
	index      int
	client     *client.Client
	downloaded int
	requested  int
	backlog    int
	buf        []byte
}

func (state *pieceProgress) readMessage() error {
	msg, err := state.client.Read()
	if err != nil {
		return err
	}
	// keep alive connection
	if msg == nil {
		return nil
	}
	switch msg.ID {
	case message.MsgUnchoke:
		state.client.Choked = false
	case message.MsgChoke:
		state.client.Choked = true
	case message.MsgHave:
		index, err := message.ParseHave(msg)
		if err != nil {
			return err
		}
		state.client.Bitfield.SetPiece(index)
	case message.MsgPiece:
		n, err := message.ParsePiece(state.index, state.buf, msg)
		if err != nil {
			return err
		}
		state.downloaded += n
		state.backlog--
	}
	return nil
}

func attemptDownloadPiece(c *client.Client, pw *pieceWork) ([]byte, error) {
	s := pieceProgress{
		index:  pw.index,
		client: c,
		buf:    make([]byte, pw.length),
	}
	c.Conn.SetDeadline(time.Now().Add(30 * time.Second))
	defer c.Conn.SetDeadline(time.Time{}) // Disable the deadline
	for s.downloaded < pw.length {
		if !s.client.Choked {
			for s.backlog < maxBacklog && s.requested < pw.length {
				blockSize := maxBlockSize
				if pw.length-s.requested < blockSize {
					blockSize = pw.length - s.requested
				}
				err := c.SendRequest(pw.index, s.requested, blockSize)
				if err != nil {
					return nil, err
				}
				s.backlog++
				s.requested += blockSize
			}
		}

		err := s.readMessage()
		if err != nil {
			return nil, err
		}
	}
	return s.buf, nil
}

func checkIntegrity(pw *pieceWork, data []byte) bool {
	hash := sha1.Sum(data)
	return bytes.Equal(hash[:], pw.hash[:])
}

func (t *Torrent) downloadTorrentWorker(peer peers.Peer, workQueue chan *pieceWork, results chan *pieceResult) {
	c, err := client.New(peer, t.InfoHash, t.PeerID)
	if err != nil {
		logger.Printf("Could not able to handshake with %s. Disconnecting...\n", peer.IP)
		return
	}
	defer c.Conn.Close()
	logger.Printf("Handshake with %s successful", peer.IP)

	c.SendUnchoke()
	c.SendInterested()

	for pw := range workQueue {
		if !c.Bitfield.HasPiece(pw.index) {
			workQueue <- pw
			continue
		}

		buf, err := attemptDownloadPiece(c, pw)
		if err != nil {
			logger.Println("Error downloading piece", pw.index, "from", peer.IP, ":", err)
			workQueue <- pw
			return
		}

		if !checkIntegrity(pw, buf) {
			logger.Println("Piece failed integrity check", pw.index, "from", peer.IP)
			workQueue <- pw
			continue
		}
		c.SendHave(pw.index)
		results <- &pieceResult{pw.index, buf}
	}
}

func (t *Torrent) calculateBoundsForPiece(index int) (begin int, end int) {
	begin = index * t.PieceLength
	end = begin + t.PieceLength
	if end > t.Length {
		end = t.Length
	}
	return begin, end
}

func (t *Torrent) calculatePieceSize(index int) int {
	begin, end := t.calculateBoundsForPiece(index)
	return end - begin
}

func (t *Torrent) Download() []byte {
	logger.Println("Starting download for", t.Name)
	workQueue := make(chan *pieceWork, len(t.PieceHashes))
	results := make(chan *pieceResult)
	for index, hash := range t.PieceHashes {
		length := t.calculatePieceSize(index)
		workQueue <- &pieceWork{index, hash, length}
	}

	for _, peer := range t.Peers {
		go t.downloadTorrentWorker(peer, workQueue, results)
	}

	buf := make([]byte, t.Length)
	bar := progressbar.NewOptions(t.Length,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetElapsedTime(false),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionSetDescription("["+t.Name+"]"),
		progressbar.OptionShowElapsedTimeOnFinish(),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))
	donePieces := 0
	for donePieces < len(t.PieceHashes) {
		res := <-results
		begin, end := t.calculateBoundsForPiece(res.index)
		copy(buf[begin:end], res.buf)
		donePieces++
		bar.Add(end - begin)
		percent := float64(donePieces) / float64(len(t.PieceHashes)) * 100
		numWorkers := runtime.NumGoroutine() - 1 // subtract 1 for main thread
		logger.Printf("(%0.2f%%) Downloaded piece #%d from %d peers\n", percent, res.index, numWorkers)
	}
	close(workQueue)
	fmt.Println()
	return buf
}
