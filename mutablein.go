package mutablein

import (
	"bytes"
	"os/exec"
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

func (m *mutableIn) Init() {
	m.isRunning = true
	exec.Command("stty", "-f", "/dev/tty", "cbreak", "min", "1").Run()
	exec.Command("stty", "-f", "/dev/tty", "-echo").Run()
}

func (m *mutableIn) Close() {
	m.isRunning = false
	exec.Command("stty", "-f", "/dev/tty", "-cbreak").Run()
	exec.Command("stty", "-f", "/dev/tty", "echo").Run()
}

func (m *mutableIn) Read(p []byte) (n int, err error) {
	if !m.isRunning {
		panic(notInitError)
	}
	n, err = m.buffer.Read(p)
	return n, err
}

func (m *mutableIn) Write(p []byte) (n int, err error) {
	if !m.isRunning {
		panic(notInitError)
	}
	n, err = m.buffer.Write(p)
	return n, err
}
