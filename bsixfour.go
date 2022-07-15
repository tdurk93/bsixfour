package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	const bitMask uint8 = 1<<6 - 1 // 0b00111111
	lookupTable := buildBase64LookupTable()

	for {
		// inBuf needs to be declared in the loop to get zeroed out
		// so it doesn't potentially corrupt the last characters if
		// the input length isn't a multiple of 3 bytes
		var inBuf [3]byte
		numBytes, err := os.Stdin.Read(inBuf[:]) // read in 3 bytes at a time

		if err != nil {
			if errors.Is(err, io.EOF) { // happens at the end of the input
				fmt.Println()
				break
			} else { // unexpected
				fmt.Fprintf(os.Stderr, "Received error when reading input: %s", err)
			}
		}

		// construct 24-bit chunk that we can dice into 4 sextets
		var threeByteChunk uint32 = uint32(inBuf[0])<<16 | uint32(inBuf[1])<<8 | uint32(inBuf[2])
		var outBuf [4]string
		for n := 0; n < len(outBuf); n++ {
			// starting with most-significant/leftmost sextet, use bitmask to extract sextets & lookup table to get ascii character
			outBuf[n] = lookupTable[uint8(threeByteChunk>>(18-n*6))&bitMask]
		}

		// at the end of the input, we might have to add padding
		for n := numBytes; n < len(inBuf); n++ {
			outBuf[n+1] = "="
		}

		// regardless of input length, output will always have a length of a multiple of 4
		fmt.Print(strings.Join(outBuf[:], ""))

	}
}

func buildBase64LookupTable() map[uint8]string {
	const capitalAOffset = 65
	const lowercaseAOffset = capitalAOffset | 1<<5 // (97)
	const zeroCharOffset = 48
	lookupTable := make(map[uint8]string)
	for n := 0; n < 26; n++ {
		lookupTable[uint8(n)] = string(rune(n + capitalAOffset))
	}
	for n := 26; n < 52; n++ {
		lookupTable[uint8(n)] = string(rune(n + lowercaseAOffset - 26))
	}
	for n := 52; n < 62; n++ {
		lookupTable[uint8(n)] = string(rune(n + zeroCharOffset - 52))
	}
	lookupTable[62] = "+"
	lookupTable[63] = "/"
	return lookupTable
}
