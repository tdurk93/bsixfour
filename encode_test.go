package main

import (
	"bytes"
	"encoding/base64"
	"math/rand"
	"testing"
)

const largeSampleDataSize = 1024 * 1024 * 1024 // 1 GB
var largeSampleData = make([]byte, largeSampleDataSize)

func checkEncodeAgainstStandardLib(t *testing.T, input []byte) {
	encodedText := Encode(bytes.NewReader(input))
	standardLibEncodedText := base64.StdEncoding.EncodeToString(input)
	if encodedText != standardLibEncodedText {
		t.Errorf("Received:\t%v\nExpected:\t%v", encodedText, standardLibEncodedText)
	}
}

func TestEncodeWithEmptyString(t *testing.T) {
	checkEncodeAgainstStandardLib(t, []byte(""))
}

func TestEncodePaddingWith1ByteInput(t *testing.T) {
	checkEncodeAgainstStandardLib(t, []byte("a"))
}

func TestEncodePaddingWith2ByteInput(t *testing.T) {
	checkEncodeAgainstStandardLib(t, []byte("ab"))
}

func TestEncodePaddingWith3ByteInput(t *testing.T) {
	checkEncodeAgainstStandardLib(t, []byte("abc"))
}

func TestEncodeWithAsciiInput(t *testing.T) {
	const base64Input = "This is sample input\nin the ascii subset of UTF-8\nwith newlines"
	checkEncodeAgainstStandardLib(t, []byte(base64Input))
}

func TestEncodeWithEmoji(t *testing.T) {
	checkEncodeAgainstStandardLib(t, []byte("👋"))
}

func TestEncodeWithNilByte(t *testing.T) {
	checkEncodeAgainstStandardLib(t, []byte{0})
}

func TestEncodeWithEveryByteValue(t *testing.T) {
	const inputLength = 256
	base64Input := make([]byte, inputLength)
	for n := 0; n < inputLength; n++ {
		base64Input[n] = byte(n)
	}
	checkEncodeAgainstStandardLib(t, base64Input)
}

func BenchmarkEncodeWithSmallInput(b *testing.B) {
	const inputLength = 256
	base64Input := make([]byte, inputLength)
	randomOffset := byte(rand.Uint32())

	for n := 0; n < inputLength; n++ {
		base64Input[n] = (byte(n) + randomOffset) % byte(inputLength-1)
	}

	reader := new(bytes.Reader)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		reader.Reset(base64Input)
		Encode(reader)
	}
}

func createLargeSampleDataSlice() []byte {
	// it's possible that largeSampleData[0] == 0 even after populating (based on the value of randomOffset),
	// but comparing largeSampleData[0] to largeSampleData[1] will correctly determine if the slice is populated
	if largeSampleData[0] != largeSampleData[1] {
		return largeSampleData
	}
	randomOffset := rand.Int()

	for n := 0; n < largeSampleDataSize; n++ {
		largeSampleData[n] = byte(n + randomOffset)
	}
	return largeSampleData
}

func BenchmarkEncodeWithLargeInput(b *testing.B) {
	base64Input := createLargeSampleDataSlice()

	reader := new(bytes.Reader)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		reader.Reset(base64Input)
		Encode(reader)
	}
}
