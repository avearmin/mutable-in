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
	return n, err
}

func (m *MutableIn) WriteByte(c byte) error {
	if !m.isRunning {
		panic(notInitError)
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	err := m.buffer.WriteByte(c)
	return err
}

func (m *MutableIn) simulateInput() {
	stdin := bufio.NewReader(os.Stdin)
	for {
		if !m.isRunning {
			return
		}
		key, err := stdin.ReadByte()
		if err != nil {
			panic(err)
		}
		err = m.WriteByte(key)
		if err != nil {
			panic(err)
		}
		if key == '\n' {
			m.cond.Signal() // We signal to waiting readers the buffer has something
			continue
		}
		fmt.Print(string(key))
	}
}
