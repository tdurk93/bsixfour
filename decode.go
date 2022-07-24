package main

import (
	"errors"
	"io"
)

func decode(reader io.Reader, dataChannel chan<- []byte) {
	decodeLookupTable := buildBase64DecoderLookupTable()

	for {
		var inBuf [4]byte
		numBytes, err := reader.Read(inBuf[:]) // read in 4 characters/bytes at a time

		if err != nil {
			if errors.Is(err, io.EOF) {
				if numBytes != 0 {
					// Input length was not a multiple of 4.
					// This tool only supports base64 with padding.
					panic("Input length not valid")
				}
				break
			} else { // unexpected
				panic(err)
			}
		}

		numUsedBytes := 3
		// construct 24-bit chunk that we can break into 3 bytes
		var fourCharacterChunk uint32 = 0
		for n, base64Char := range inBuf {
			// first perform padding check so we can avoid using the lookup table
			// on padding characters where there isn't technically an output mapping
			if n >= 2 && rune(base64Char) == '=' { // only check 3rd & 4th characters for padding
				numUsedBytes = n - 1
				break
			}
			fourCharacterChunk = fourCharacterChunk | uint32(decodeLookupTable[rune(inBuf[n])])<<(6*(3-n))
		}
		outBuf := [3]byte{byte(fourCharacterChunk >> 16), byte(fourCharacterChunk >> 8), byte(fourCharacterChunk)}
		dataChannel <- outBuf[0:numUsedBytes]
	}
	close(dataChannel)
}

func Decode(reader io.Reader) []byte {
	data := make([]byte, 200) // arbitrarily choosing an initial capacity of 200
	originalDataChannel := make(chan []byte)
	go decode(reader, originalDataChannel)
	for {
		val, isOpen := <-originalDataChannel
		if !isOpen {
			break
		}
		data = append(data, val...)
	}
	return data
}

func buildBase64DecoderLookupTable() map[rune]byte {
	lookupTable := make(map[rune]byte)
	for n := 0; n < 26; n++ {
		lookupTable[rune(n+capitalAOffset)] = byte(n)
	}
	for n := 26; n < 52; n++ {
		lookupTable[rune(n+lowercaseAOffset-26)] = byte(n)
	}
	for n := 52; n < 62; n++ {
		lookupTable[rune(n+zeroCharOffset-52)] = byte(n)
	}
	lookupTable['+'] = 62
	lookupTable['/'] = 63
	return lookupTable
}
