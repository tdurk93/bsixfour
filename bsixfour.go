package main

import (
	"fmt"
	"os"
)

func main() {
	// TODO add ability to either encode or decode from command line
	// base64Channel := make(chan string)
	// go encode(os.Stdin, base64Channel)
	// for {
	// 	val, isOpen := <-base64Channel
	// 	if !isOpen {
	// 		break
	// 	}
	// 	fmt.Print(val)
	// }

	originalDataChannel := make(chan []byte)
	go decode(os.Stdin, originalDataChannel)
	for {
		val, isOpen := <-originalDataChannel
		if !isOpen {
			break
		}
		fmt.Print(string(val))
	}

	// printing final newline using stderr
	// so it doesn't corrupt output when piped to other applications.
	// Eventually might include a flag to make this behavior optional
	fmt.Fprintf(os.Stderr, "\n")
}
