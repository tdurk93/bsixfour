package main

import (
	"errors"
	"io"
	"strings"
)

func encode(reader io.Reader, base64Channel chan<- string) {
	lookupTable := buildBase64EncoderLookupTable()

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
			// starting with most-significant/leftmost sextet,
			// use bitmask to extract sextets & lookup table to get ascii character
			outBuf[n] = lookupTable[byte(threeByteChunk>>(18-n*6))&sextetBitMask]
		}

		// at the end of the input, we might have to add padding
		for n := numBytes; n < len(inBuf); n++ {
			outBuf[n+1] = '='
		}

		// regardless of input length, output will always have a length of a multiple of 4
		base64Channel <- string(outBuf[:])
	}
	close(base64Channel)
}

func Encode(reader io.Reader) string {
	stringBuilder := new(strings.Builder)
	base64DataChannel := make(chan string)
	go encode(reader, base64DataChannel)
	for {
		val, isOpen := <-base64DataChannel
		if !isOpen {
			break
		}
		stringBuilder.Write([]byte(val))
	}
	return stringBuilder.String()
}

func buildBase64EncoderLookupTable() map[byte]rune {
	lookupTable := make(map[byte]rune)
	for n := 0; n < 26; n++ {
		lookupTable[byte(n)] = rune(n + capitalAOffset)
	}
	for n := 26; n < 52; n++ {
		lookupTable[byte(n)] = rune(n + lowercaseAOffset - 26)
	}
	for n := 52; n < 62; n++ {
		lookupTable[byte(n)] = rune(n + zeroCharOffset - 52)
	}
	lookupTable[62] = '+'
	lookupTable[63] = '/'
	return lookupTable
}
