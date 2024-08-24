package torrentfile

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/Harry-kp/nebula/peers"
	"github.com/jackpal/bencode-go"
)

type bencodeTrackerResp struct {
	Interval int    `bencode:"interval"`
	Peers    string `bencode:"peers"`
}

func (tf *TorrentFile) createTrackerURL(peer_id [20]byte, port uint16) (string, error) {
	baseURL, err := url.Parse(tf.Announce)
	if err != nil {
		return "", err
	}

	params := url.Values{}

	params.Add("peer_id", string(peer_id[:]))
	params.Add("info_hash", string(tf.InfoHash[:]))
	params.Add("port", strconv.Itoa(int(port)))
	params.Add("uploaded", "0")
	params.Add("downloaded", "0")
	params.Add("compact", "1")
	params.Add("left", strconv.Itoa(tf.Length))
	return fmt.Sprintf("%s?%s", baseURL, params.Encode()), nil
}
func (tf *TorrentFile) fetchPeers(peer_id [20]byte, port uint16) ([]peers.Peer, error) {
	if len(tf.Announce) < 7 {
		return nil, fmt.Errorf("Invalid announce URL")
	}

	switch tf.Announce[:7] {
	case "http://":
		return tf.fetchPeersHttp(peer_id, port)
	case "udp://":
		return nil, fmt.Errorf("UDP tracker not supported yet.We are working on it")
	default:
		return nil, fmt.Errorf("Currently, we only support HTTP trackers protocol. Please use torrent with http:// announce URL. We are working on it")
	}
}

func (tf *TorrentFile) fetchPeersHttp(peer_id [20]byte, port uint16) ([]peers.Peer, error) {
	trackerURL, err := tf.createTrackerURL(peer_id, port)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 15 * time.Second}
	response, err := client.Get(trackerURL)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	trackerResp := bencodeTrackerResp{}
	err = bencode.Unmarshal(response.Body, &trackerResp)
	if err != nil {
		return nil, err
	}
	return peers.Unmarshal([]byte(trackerResp.Peers))
}
