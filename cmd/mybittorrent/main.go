package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	// bencode "github.com/jackpal/bencode-go" // Available if you need it!
)

func decodeBencode(bencodedString string, start int) (interface{}, int, error) {
	if start >= len(bencodedString) {
		return "", start, fmt.Errorf("Invalid Bencoded string")
	}

	if len(bencodedString) == 0 {
		return "", start, fmt.Errorf("Invalid Bencoded string")
	}

	if bencodedString[start] == 'l' {
		return decodeList(bencodedString, start)
	} else if bencodedString[start] == 'i' {
		return decodeInt(bencodedString, start)
	} else if bencodedString[start] >= '0' && bencodedString[start] <= '9' {
		return decodeString(bencodedString, start)
	} else {
		return "", start, fmt.Errorf("Invalid Bencoded string")
	}
}

func decodeList(bencodedString string, start int) (interface{}, int, error) {
	start++
	result := make([]interface{}, 0)

	for start < len(bencodedString) && bencodedString[start] != 'e' {
		var decoded interface{}
		var err error

		decoded, start, err = decodeBencode(bencodedString, start)
		if err != nil {
			return "", start, err
		}

		result = append(result, decoded)
	}
	return result, start + 1, nil
}

func decodeInt(bencodedString string, start int) (interface{}, int, error) {
	start++
	end := start

	for end < len(bencodedString) && bencodedString[end] != 'e' {
		end++
	}

	if end == len(bencodedString) {
		return "", start, fmt.Errorf("Invalid Bencoded string")
	}

	value, err := strconv.Atoi(bencodedString[start:end])
	if err != nil {
		return "", start, err
	}

	return value, end + 1, nil
}

func decodeString(bencodedString string, start int) (interface{}, int, error) {
	colonIndex := start

	for colonIndex < len(bencodedString) && bencodedString[colonIndex] != ':' {
		colonIndex++
	}

	if colonIndex == len(bencodedString) {
		return "", start, fmt.Errorf("Invalid Bencoded string")
	}

	length, err := strconv.Atoi(bencodedString[start:colonIndex])
	if err != nil {
		return "", start, err
	}

	start = colonIndex + 1
	end := start + length

	if end > len(bencodedString) {
		return "", start, fmt.Errorf("Invalid Bencoded string")
	}

	return bencodedString[start:end], end, nil
}

func main() {
	command := os.Args[1]

	if command == "decode" {
		bencodedValue := os.Args[2]

		decoded, _, err := decodeBencode(bencodedValue, 0)
		if err != nil {
			fmt.Println(err)
			return
		}

		jsonOutput, _ := json.Marshal(decoded)
		fmt.Println(string(jsonOutput))
	} else {
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}
