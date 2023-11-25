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
	m.buffer.cond.Signal() // We signal to waiting readers the buffer has something
	m.cursor = -1
	m.Write([]byte{'\n'})
}

func handleBackspace(m *MutableIn) {
	if m.cursor < 1 {
		return
	}

	if m.cursor >= m.buffer.Len() {
		m.buffer.truncate()
		m.cursor--
		fmt.Print("\b \b")
	} else {
		m.buffer.deleteIndexAt(m.cursor - 1)
		m.cursor--
		fmt.Print("\033[D")
		redrawLineFromCursor(m)
		fmt.Print("\033[K")
		returnCursorToPos(m)
	}
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

// Does  not appear in key events map. Consider this the default
func handleOtherKeys(m *MutableIn, key [3]byte, n int) {
	if m.cursor >= m.buffer.Len() {
		m.Write(key[:n])
	} else {
		m.buffer.insert(key[:n], m.cursor)
		redrawLineFromCursor(m)
		m.cursor++
		returnCursorToPos(m)
	}

}

// helpers
func redrawLineFromCursor(m *MutableIn) {
	lineAfterCursor := m.buffer.stringFromIndex(m.cursor)
	fmt.Print(lineAfterCursor)
}

func returnCursorToPos(m *MutableIn) {
	pos := m.buffer.Len() - m.cursor
	fmt.Printf("\033[%dD", pos)
}
