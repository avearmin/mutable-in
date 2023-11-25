package mutablein

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

type MutableIn struct {
	buffer    *buffer
	isRunning bool
	cursor    int
}

func NewMutableIn() *MutableIn {
	buf := newbuffer()
	muIn := MutableIn{
		buffer: buf,
	}
	return &muIn
}

func (m *MutableIn) Init() {
	m.isRunning = true
	exec.Command("stty", "-f", "/dev/tty", "cbreak", "min", "1").Run()
	exec.Command("stty", "-f", "/dev/tty", "-echo").Run()
	go m.simulateInput()
}

func (m *MutableIn) Close() {
	m.isRunning = false
	exec.Command("stty", "-f", "/dev/tty", "-cbreak").Run()
	exec.Command("stty", "-f", "/dev/tty", "echo").Run()
}

func (m *MutableIn) Read(p []byte) (n int, err error) {
	if !m.isRunning {
		return 0, notInitError
	}
	n, err = m.buffer.Read(p)
	return n, err
}

func (m *MutableIn) Write(p []byte) (n int, err error) {
	if !m.isRunning {
		return 0, notInitError
	}
	n, err = m.buffer.Write(p)
	m.cursor += n
	fmt.Print(string(p))
	return n, nil
}

func (m *MutableIn) simulateInput() {
	stdin := bufio.NewReader(os.Stdin)
	keyEvents := keyEvents()
	for {
		if !m.isRunning {
			return
		}
		var key [3]byte
		n, err := stdin.Read(key[:])
		if err != nil {
			panic(err)
		}
		event, hasEvent := keyEvents[key]
		if hasEvent {
			event.callback(m)
			continue
		}
		handleOtherKeys(m, key, n)
	}
}
