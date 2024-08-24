package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Harry-kp/nebula/logger"
	"github.com/Harry-kp/nebula/torrentfile"
	"github.com/Harry-kp/nebula/utils"
)

// resolveFilePath resolves the absolute path and checks if the file exists
func resolveFilePath(path string) (string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	return absPath, nil
}

func main() {
	// Define flags for input and output file paths
	inputFile := flag.String("input", "", "Path to the input torrent file (required)")
	outputFile := flag.String("output", ".", "Path to the output file or directory (default: current directory)")
	logEnabled := flag.Bool("log", false, "Enable logging")

	// Parse the flags
	flag.Parse()

	// Validate required flags
	if *inputFile == "" {
		fmt.Println("Torrent file path is required")
		flag.Usage()
		os.Exit(1)
	}

	// Initialize the global logger
	config := logger.Config{
		LogEnabled: *logEnabled,
	}
	logger.Init(config)

	// Display banner
	utils.Banner()

	// Resolve input file path
	inPath, err := resolveFilePath(*inputFile)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Error resolving input file path: %v", err))
	}

	// Resolve output file path
	outPath, err := resolveFilePath(*outputFile)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Error resolving output file path: %v", err))
	}

	// Open the torrent file
	tf, err := torrentfile.Open(inPath)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Error opening torrent file: %v", err))
	}

	// If output path is a directory, append the torrent file name
	if stat, err := os.Stat(outPath); err == nil && stat.IsDir() {
		outPath = filepath.Join(outPath, tf.Name)
	}

	// Validate the output path
	if _, err := os.Stat(outPath); err == nil {
		logger.Fatal("Error: output file already exists")
	}

	// Download the torrent file to the specified output path
	err = tf.DownloadToFile(outPath)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Error downloading torrent file: %v", err))
	}

	fmt.Println("Download completed successfully")
}
