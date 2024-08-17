# Nebula: A Go-Based Torrent Client 🚀

Nebula is a command-line torrent client written in Go. It's designed to be a lightweight and efficient tool for downloading torrent files.

### Credits where it's due:

- I got the insipiration to build this project from [here](https://blog.jse.li/posts/torrent/) and this source code is almost similar to [this](https://github.com/veggiedefender/torrent-client).
- I just write this in my repostory to learn more about the BitTorrent protocol and build a functional torrent client from scratch. It helped me to understand the working of torrent clients and how they communicate with trackers and peers to download files.
- I will try to add more features to this project and make it more efficient and user-friendly.

### Current Features:

- **Parsing .torrent files:** Nebula can parse .torrent files and extract relevant information such as the announce URL, file list, and piece hashes.
- **Downloading torrent content:** Nebula can download the content of a torrent file using the information extracted from the .torrent file. ⬇️
- **HTTP Tracker Support:** Nebula can communicate with HTTP trackers to find peers for downloading torrent content.
- **Piece Management:** Nebula efficiently manages pieces of the torrent, downloading them concurrently from multiple peers.
- **Data Verification:** Nebula verifies the integrity of downloaded pieces using SHA-1 hashes to ensure data accuracy. ✅

### Future Features (Planned):

- **Magnet Link Support:** Add support for downloading torrents using magnet links. 🧲
- **UDP Tracker Support:** Implement compatibility with UDP trackers for peer discovery.
- **Endgame Mode:** Optimize the download process in the final stages to ensure all pieces are acquired.
- **Selective Downloading:** Allow users to choose specific files to download from a torrent.
- **Prioritization:** Enable prioritization of specific files or pieces for faster access.
- **Sequential Downloading:** Implement an option for downloading files sequentially for smoother playback of media files.
- **Configuration Options:** Introduce a configuration file or command-line flags for customization.
- **Improved User Interface:** Enhance the command-line interface or potentially explore a graphical user interface (GUI) for a more user-friendly experience.

### Getting Started:

1. **Prerequisites:**

   - Go (version 1.16 or later) installed on your system.

2. **Installation:**

   **Option 1: Download Pre-built Binaries**

   - Download the appropriate binary for your operating system from the [Releases](https://github.com/Harry-kp/nebula/releases) page.
   - Make the binary executable (e.g., `chmod +x nebula-linux`).
   - Move the binary to a directory in your PATH (e.g., `/usr/local/bin`).

   **Option 2: Build from Source**

   ```bash
   go get github.com/Harry-kp/nebula
   ```

3. **Usage:**

   ```bash
   nebula <path/to/torrent.torrent> <path/to/output>
   ```

**Example:**

```bash
nebula my_favorite_movie.torrent .
```

### How it Works:

1. **Parsing:** Nebula parses the `.torrent` file to extract essential information, including the announce URL, file list, piece hashes, and total size.
2. **Tracker Communication:** Nebula contacts the tracker specified in the announce URL to obtain a list of peers participating in the torrent swarm.
3. **Peer Connection:** Nebula establishes connections with multiple peers from the list provided by the tracker.
4. **Piece Downloading:** Nebula requests pieces of the torrent from different peers, prioritizing pieces that are rare among the connected peers.
5. **Data Verification:** As pieces are downloaded, Nebula verifies their integrity using the SHA-1 hashes included in the `.torrent` file.
6. **Assembly:** Once all the pieces are downloaded and verified, Nebula assembles them into the complete files as specified in the `.torrent` file.

### Contributing:

Contributions are welcome! If you'd like to contribute to Nebula, please follow these steps:

1. Fork the repository.
2. Create a new branch for your feature or bug fix.
3. Make your changes and commit them with descriptive commit messages.
4. Push your changes to your fork.
5. Submit a pull request to the main repository.

### License:

Nebula is licensed under the MIT License. See the `LICENSE` file for details.

### Acknowledgments:

- This project was inspired by the desire to learn more about the BitTorrent protocol and build a functional torrent client from scratch.
- Thanks to the Go community for providing excellent resources and libraries.

### Disclaimer:

This project is for educational purposes and personal use. Please respect copyright laws and download content responsibly.
