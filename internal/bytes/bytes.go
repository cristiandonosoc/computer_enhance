// Package bytes holds utilities for working with bytes.
package bytes

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


