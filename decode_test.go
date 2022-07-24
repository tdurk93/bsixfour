package main

import (
	"bytes"
	"encoding/base64"
	"math"
	"math/rand"
	"strings"
	"testing"
)

var largeBase64String = ""

func checkDecodeAgainstStandardLib(t *testing.T, expectedDecodedData []byte) {
	// it's simpler to pass in the decoded data and use
	// base64 library's EncodeToString() to build valid input to Decode()
	input := base64.StdEncoding.EncodeToString(expectedDecodedData)
	decodedData := Decode(strings.NewReader(input))
	if bytes.Equal(decodedData, expectedDecodedData) {
		t.Errorf("Received:\t%v\nExpected:\t%v", decodedData, expectedDecodedData)
	}
}

func TestDecodeWithEmptyString(t *testing.T) {
	checkDecodeAgainstStandardLib(t, []byte(""))
}

func TestDecodePaddingWith1ByteOutput(t *testing.T) {
	checkDecodeAgainstStandardLib(t, []byte("a"))
}

func TestDecodePaddingWith2ByteOutput(t *testing.T) {
	checkDecodeAgainstStandardLib(t, []byte("ab"))
}

func TestDecodePaddingWith3ByteOutput(t *testing.T) {
	checkDecodeAgainstStandardLib(t, []byte("abc"))
}

func TestDecodeWithAsciiInput(t *testing.T) {
	const base64Input = "This is sample input\nin the ascii subset of UTF-8\nwith newlines"
	checkDecodeAgainstStandardLib(t, []byte(base64Input))
}

func TestDecodeWithEmoji(t *testing.T) {
	checkDecodeAgainstStandardLib(t, []byte("👋"))
}

func TestDecodeWithNilByte(t *testing.T) {
	checkDecodeAgainstStandardLib(t, []byte{0})
}

func TestDecodeWithEveryByteValue(t *testing.T) {
	const inputLength = 256
	decodedData := make([]byte, inputLength)
	for n := range decodedData {
		decodedData[n] = byte(n)
	}
	checkDecodeAgainstStandardLib(t, decodedData)
}

func createBase64StringOfSmallInput() string {
	const decodedDataLength = 256
	decodedData := make([]byte, decodedDataLength)
	randomOffset := byte(rand.Uint32())

	for n := range decodedData {
		decodedData[n] = (byte(n) + randomOffset) % byte(decodedDataLength-1)
	}
	return base64.StdEncoding.EncodeToString(decodedData)
}

func BenchmarkDecodeSmallInput(b *testing.B) {
	base64String := createBase64StringOfSmallInput()
	reader := new(strings.Reader)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		reader.Reset(base64String)
		Decode(reader)
	}
}

func BenchmarkDecodeSmallInputWithStandardLib(b *testing.B) {
	base64String := createBase64StringOfSmallInput()
	unusedReader := new(strings.Reader)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		unusedReader.Reset(base64String) // included for fair comparison with my Decode()
		base64.StdEncoding.DecodeString(base64String)
	}
}

func createBase64StringOfLargeInput() string {
	// in order to avoid allocating 1GB of data plus a 1.33 GB string for its encoding,
	// we build up the encoded string in smaller chunks so the garbage collector can
	// clean up the original data as we go

	if len(largeBase64String) > 0 {
		return largeBase64String // use memoizing for subsequent calls
	}

	const oneMB = 1024 * 1024
	const smallBufferSize = oneMB * 3  // value should be divisble by 3 to avoid inserting padding prematurely
	const bigBufferSize = oneMB * 1024 // 1 GB
	decodedData := make([]byte, smallBufferSize)
	base64StringBuffer := make([]byte, int(math.Ceil(float64(bigBufferSize)/3))*4)
	randomOffset := rand.Int()
	currIndex := 0
	var endIndex int
	for currIndex < bigBufferSize {
		endIndex = currIndex + smallBufferSize
		// on the last loop, don't iterate as many times
		if endIndex > bigBufferSize {
			// decodedData = decodedData[0 : endIndex%bigBufferSize]
			decodedData = decodedData[0 : bigBufferSize-currIndex]
			endIndex = bigBufferSize
		}

		// since currIndex should always be divisible by 3,
		// this should yield an integer without truncating/casting
		bigBufferSlice := base64StringBuffer[currIndex*4/3:]

		for ; currIndex < endIndex; currIndex++ {
			decodedData[currIndex%smallBufferSize] = byte(currIndex + randomOffset)
		}

		base64.StdEncoding.Encode(bigBufferSlice, decodedData)
	}
	largeBase64String = string(base64StringBuffer)

	return largeBase64String
}

func BenchmarkDecodeLargeInput(b *testing.B) {
	base64String := []byte(createBase64StringOfLargeInput())

	reader := new(bytes.Reader)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		reader.Reset(base64String)
		Decode(reader)
	}
}

func BenchmarkDecodeLargeInputWithStandardLib(b *testing.B) {
	base64String := createBase64StringOfLargeInput()

	reader := new(strings.Reader)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		reader.Reset(base64String)
		base64.StdEncoding.DecodeString(base64String)
	}
}
