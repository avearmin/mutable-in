package mutablein

import "fmt"

var (
	UpArrow    = [3]byte{0x1b, 0x5b, 0x41}
	DownArrow  = [3]byte{0x1b, 0x5b, 0x42}
	RightArrow = [3]byte{0x1b, 0x5b, 0x43}
	LeftArrow  = [3]byte{0x1b, 0x5b, 0x44}
	Enter      = [3]byte{0x0a}
)

type event struct {
	callback func(m *MutableIn)
}

func keyEvents() map[[3]byte]event {
	return map[[3]byte]event{
		Enter:      {callback: handleEnter},
		RightArrow: {callback: handleRightArrow},
		LeftArrow:  {callback: handleLeftArrow},
		UpArrow:    {callback: handleUpDownArrows},
		DownArrow:  {callback: handleUpDownArrows},
	}
}

func handleEnter(m *MutableIn) {
	m.cond.Signal() // We signal to waiting readers the buffer has something
	m.cursor = 0
	m.Write([]byte{'\n'})
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
