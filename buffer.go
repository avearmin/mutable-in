package mutablein

import (
	"io"
	"sync"
)

type buffer struct {
	bytes []byte
	mu    *sync.Mutex
	cond  *sync.Cond
}

func newbuffer() *buffer {
	buf := buffer{
		bytes: []byte{},
		mu:    &sync.Mutex{},
	}
	buf.cond = sync.NewCond(buf.mu)
	return &buf
}

func (buf *buffer) Read(p []byte) (n int, err error) {
	buf.mu.Lock()
	defer buf.mu.Unlock()

	buf.cond.Wait() // We want readers to wait until the enter key is pressed

	if buf.Len() == 0 {
		return 0, io.EOF
	}
	n = copy(p, buf.bytes)
	buf.bytes = buf.bytes[n:]

	return n, nil
}

func (buf *buffer) Write(p []byte) (n int, err error) {
	buf.mu.Lock()
	defer buf.mu.Unlock()

	buf.bytes = append(buf.bytes, p...)
	n = len(p)
	return n, nil
}

func (buf *buffer) insert(p []byte, pos int) {
	buf.mu.Lock()
	defer buf.mu.Unlock()

	data := buf.Bytes()
	first := data[:pos]
	second := data[pos:]

	newData := make([]byte, 0, len(first)+len(p)+len(second))
	newData = append(newData, first...)
	newData = append(newData, p...)
	newData = append(newData, second...)

	buf.bytes = newData
}

func (buf *buffer) deleteIndexAt(pos int) {
	buf.mu.Lock()
	defer buf.mu.Unlock()

	data := buf.Bytes()
	first := data[:pos]
	second := data[pos+1:]

	newData := append(first, second...)
	buf.bytes = newData
}

func (buf *buffer) truncate() {
	buf.mu.Lock()
	defer buf.mu.Unlock()

	buf.bytes = buf.bytes[:buf.Len()-1]
}

func (buf *buffer) Len() int {
	return len(buf.bytes)
}

func (buf *buffer) String() string {
	return string(buf.bytes)
}

func (buf *buffer) stringFromIndex(index int) string {
	return string(buf.bytes[index:])
}

func (buf *buffer) Bytes() []byte {
	return buf.bytes
}
