package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

const sextetBitMask byte = 1<<6 - 1 // 0b00111111

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
	const capitalAOffset = 65
	const lowercaseAOffset = capitalAOffset | 1<<5 // (97)
	const zeroCharOffset = 48
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
		for n := 0; n < len(inBuf); n++ {
			// first perform padding check so we can avoid using the lookup table
			// on padding characters where there isn't technically an output mapping
			if n >= 2 && rune(inBuf[n]) == '=' { // only check 3rd & 4th characters for padding
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

func Decode(reader io.Reader) string {
	stringBuilder := new(strings.Builder)
	originalDataChannel := make(chan []byte)
	go decode(reader, originalDataChannel)
	for {
		val, isOpen := <-originalDataChannel
		if !isOpen {
			break
		}
		stringBuilder.Write(val)
	}
	return stringBuilder.String()
}

func buildBase64DecoderLookupTable() map[rune]byte {
	const capitalAOffset = 65
	const lowercaseAOffset = capitalAOffset | 1<<5 // (97)
	const zeroCharOffset = 48
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
