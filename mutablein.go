package mutablein

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sync"
)

type MutableIn struct {
	buffer    *bytes.Buffer
	isRunning bool
	mu        *sync.Mutex
	cond      *sync.Cond
	cursor    int
}

func NewMutableIn() *MutableIn {
	muIn := MutableIn{
		buffer: bytes.NewBuffer([]byte{}),
		mu:     &sync.Mutex{},
	}
	muIn.cond = sync.NewCond(muIn.mu)
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
		panic(notInitError)
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	m.cond.Wait() // We want readers to wait until the enter key is pressed

	n, err = m.buffer.Read(p)
	return n, err
}

func (m *MutableIn) Write(p []byte) (n int, err error) {
	if !m.isRunning {
		panic(notInitError)
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	n, err = m.buffer.Write(p)
	m.cursor += len(p)
	fmt.Print(string(p))
	return n, err
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
		n, err = m.Write(key[:n])
		if err != nil {
			panic(err)
		}
	}
}
