package mutablein

import (
	"bytes"

	"golang.org/x/term"
)

type mutableIn struct {
	buffer    *bytes.Buffer
	isRunning bool
}

func NewMutableIn() *mutableIn {
	return &mutableIn{
		buffer: bytes.NewBuffer([]byte{}),
	}
}

func (m *mutableIn) Init() *term.State {
	m.isRunning = true
	oldState, err := term.MakeRaw(0)
	if err != nil {
		panic(err)
	}
	return oldState
}

func (m *mutableIn) Close(oldState *term.State) {
	m.isRunning = false
	err := term.Restore(0, oldState)
	if err != nil {
		panic(err)
	}
}

func (m *mutableIn) Read(p []byte) (n int, err error) {
	n, err = m.buffer.Read(p)
	return n, err
}

func (m *mutableIn) Write(p []byte) (n int, err error) {
	n, err = m.buffer.Write(p)
	return n, err
}
