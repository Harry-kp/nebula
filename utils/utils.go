package utils

import "fmt"

func Banner() {
	green := "\033[32m" // ANSI escape code for green text
	reset := "\033[0m"  // ANSI escape code to reset color

	fmt.Println(green + `    ======================================
     _   _        _             _
    | \ | |      | |           | |
    |  \| |  ___ | |__   _   _ | |  __ _
    | . ` + "`" + ` | / _ \| '_ \ | | | || | / _` + "`" + ` |
    | |\  ||  __/| |_) || |_| || || (_| |
    |_| \_| \___||_.__/  \__,_||_| \__,_|
    ======================================
` + reset)
}
