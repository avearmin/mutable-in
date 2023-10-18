package mutablein

import "bytes"

type mutableIn struct {
	buffer *bytes.Buffer
}

func NewMutableIn() *mutableIn {
	return &mutableIn{
		buffer: bytes.NewBuffer([]byte{}),
	}
}

func (m mutableIn) Read(p []byte) (n int, err error) {
	n, err = m.buffer.Read(p)
	return n, err
}

func (m mutableIn) Write(p []byte) (n int, err error) {
	n, err = m.buffer.Write(p)
	return n, err
}
