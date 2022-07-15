package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

const sextetBitMask uint8 = 1<<6 - 1 // 0b00111111

func main() {
	base64Output := Encode(os.Stdin)
	fmt.Println(base64Output)
}

func Encode(reader io.Reader) string {
	// note that this is currently not a great implementation
	// because it stores the entire contents of the base64
	// string in memory. I will improve with goroutines/channelss
	// at a future date.
	lookupTable := buildBase64EncoderLookupTable()
	stringBuilder := new(strings.Builder)

	for {
		// inBuf needs to be declared in the loop to get zeroed out
		// so it doesn't potentially corrupt the last characters if
		// the input length isn't a multiple of 3 bytes
		var inBuf [3]byte
		numBytes, err := reader.Read(inBuf[:]) // read in 3 bytes at a time

		if err != nil {
			if errors.Is(err, io.EOF) { // happens at the end of the input
				break
			} else { // unexpected
				panic(err)
			}
		}

		// construct 24-bit chunk that we can dice into 4 sextets
		var threeByteChunk uint32 = uint32(inBuf[0])<<16 | uint32(inBuf[1])<<8 | uint32(inBuf[2])
		var outBuf [4]rune
		for n := 0; n < len(outBuf); n++ {
			// starting with most-significant/leftmost sextet, use bitmask to extract sextets & lookup table to get ascii character
			outBuf[n] = lookupTable[uint8(threeByteChunk>>(18-n*6))&sextetBitMask]
		}

		// at the end of the input, we might have to add padding
		for n := numBytes; n < len(inBuf); n++ {
			outBuf[n+1] = '='
		}

		// regardless of input length, output will always have a length of a multiple of 4
		stringBuilder.WriteString(string(outBuf[:]))
	}

	return stringBuilder.String()
}

func buildBase64EncoderLookupTable() map[uint8]rune {
	const capitalAOffset = 65
	const lowercaseAOffset = capitalAOffset | 1<<5 // (97)
	const zeroCharOffset = 48
	lookupTable := make(map[uint8]rune)
	for n := 0; n < 26; n++ {
		lookupTable[uint8(n)] = rune(n + capitalAOffset)
	}
	for n := 26; n < 52; n++ {
		lookupTable[uint8(n)] = rune(n + lowercaseAOffset - 26)
	}
	for n := 52; n < 62; n++ {
		lookupTable[uint8(n)] = rune(n + zeroCharOffset - 52)
	}
	lookupTable[62] = '+'
	lookupTable[63] = '/'
	return lookupTable
}
