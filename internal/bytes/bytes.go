// Package bytes holds utilities for working with bytes.
package bytes

import (
	"fmt"
)

// ToUint16 merges the two bytes into a little-endian 16 bit value.
// b1 would be the lowest byte.
func ToUint16(b1, b2 byte) uint16 {
	return uint16(b1) | (uint16(b2) << 8)
}

// ByteStream --------------------------------------------------------------------------------------

// ByteStream is a helper data structure to treat a an array of bytes as a stream.
type ByteStream struct {
	Data []byte
	Pos  int
}

func NewByteStream(data []byte) *ByteStream {
	return &ByteStream{
		Data: data,
	}
}

func (bs *ByteStream) Peek() byte {
	return bs.Data[bs.Pos]
}

func (bs *ByteStream) Advance() byte {
	b := bs.Data[bs.Pos]
	bs.Pos++
	return b
}

func (bs *ByteStream) IsEOF() bool {
	return bs.Pos == len(bs.Data)
}

func (bs *ByteStream) ReadByte() (byte, error) {
	if bs.IsEOF() {
		return 0, fmt.Errorf("stream EOF")
	}

	b := bs.Advance()
	return b, nil
}

func (bs *ByteStream) ReadWord() (uint16, error) {
	b1, b2, err := bs.ReadWordAsBytes()
	if err != nil {
		return 0, err
	}

	return ToUint16(b1, b2), nil
}

func (bs *ByteStream) ReadWordAsBytes() (_b1, _b2 byte, _err error) {
	b1, err := bs.ReadByte()
	if err != nil {
		return 0, 0, fmt.Errorf("reading first byte: %w", err)
	}

	b2, err := bs.ReadByte()
	if err != nil {
		return 0, 0, fmt.Errorf("reading second byte: %w", err)
	}

	return b1, b2, nil
}
