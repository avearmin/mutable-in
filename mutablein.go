package mutablein

import (
	"bytes"
	"os/exec"
)

type MutableIn struct {
	buffer    *bytes.Buffer
	isRunning bool
}

func NewMutableIn() *MutableIn {
	return &MutableIn{
		buffer: bytes.NewBuffer([]byte{}),
	}
}

func (m *MutableIn) Init() {
	m.isRunning = true
	exec.Command("stty", "-f", "/dev/tty", "cbreak", "min", "1").Run()
	exec.Command("stty", "-f", "/dev/tty", "-echo").Run()
}

func (m *MutableIn) Close() {
	m.isRunning = false
	exec.Command("stty", "-f", "/dev/tty", "-cbreak").Run()
	exec.Command("stty", "-f", "/dev/tty", "echo").Run()
}

func (m *MutableIn) Read(p []byte) (n int, err error) {
	if !m.isRunning {
		panic(notInitError)
	}
	n, err = m.buffer.Read(p)
	return n, err
}

func (m *MutableIn) Write(p []byte) (n int, err error) {
	if !m.isRunning {
		panic(notInitError)
	}
	n, err = m.buffer.Write(p)
	return n, err
}
