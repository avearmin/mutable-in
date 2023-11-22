package mutablein

import "fmt"

var (
	UpArrow    = [3]byte{0x1b, 0x5b, 0x41}
	DownArrow  = [3]byte{0x1b, 0x5b, 0x42}
	RightArrow = [3]byte{0x1b, 0x5b, 0x43}
	LeftArrow  = [3]byte{0x1b, 0x5b, 0x44}
	Enter      = [3]byte{0x0a}
	Backspace  = [3]byte{0x7f}
)

type event struct {
	callback func(m *MutableIn)
}

func keyEvents() map[[3]byte]event {
	return map[[3]byte]event{
		Enter:      {callback: handleEnter},
		Backspace:  {callback: handleBackspace},
		RightArrow: {callback: handleRightArrow},
		LeftArrow:  {callback: handleLeftArrow},
		UpArrow:    {callback: handleUpDownArrows},
		DownArrow:  {callback: handleUpDownArrows},
	}
}

// handlers
func handleEnter(m *MutableIn) {
	m.cond.Signal() // We signal to waiting readers the buffer has something
	m.cursor = -1   // Don't know whats up here. If I set the cursor to 0, then backspace does not work properly. But -1 works perfectly. What am I missing? Is 0 the wrong way to think about reseting cursor pos?
	m.Write([]byte{'\n'})
}

func handleBackspace(m *MutableIn) {
	if m.cursor < 1 {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	var newData []byte
	if m.cursor >= m.buffer.Len() {
		newData = processDataWhenBackspaceAtEnd(m)
	} else {
		newData = proccessDataWhenBackspaceAtMiddle(m)
	}

	m.buffer.Reset()
	m.buffer.Write(newData)

	m.cursor--
	redrawTermLine(m)
}

func handleRightArrow(m *MutableIn) {
	if m.cursor > m.buffer.Len()-1 {
		return
	}
	m.cursor++
	fmt.Print("\033[C")

}

func handleLeftArrow(m *MutableIn) {
	if m.cursor < 1 {
		return
	}
	m.cursor--
	fmt.Print("\033[D")
}

func handleUpDownArrows(m *MutableIn) {
	// These keys should do nothing
}

// helpers
func processDataWhenBackspaceAtEnd(m *MutableIn) []byte {
	data := m.buffer.Bytes()
	newData := data[:len(data)-1]
	return newData
}

func proccessDataWhenBackspaceAtMiddle(m *MutableIn) []byte {
	data := m.buffer.Bytes()
	first := data[:m.cursor]
	second := data[m.cursor+1:]

	newData := append(first, second...)
	return newData
}

func redrawTermLine(m *MutableIn) {
	fmt.Print("\033[1G")
	fmt.Print("\033[2K")
	fmt.Print(m.buffer.String())
	fmt.Print("\r")
	for i := 0; i < m.cursor; i++ {
		fmt.Print("\033[C")
	}
}
