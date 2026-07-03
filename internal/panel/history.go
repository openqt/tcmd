package panel

// History tracks visited directories.
type History struct {
	entries []string
	index   int
}

// NewHistory creates empty history.
func NewHistory() *History {
	return &History{index: -1}
}

// Push records a path if different from current position.
func (h *History) Push(path string) {
	if h == nil {
		return
	}
	if h.index >= 0 && h.entries[h.index] == path {
		return
	}
	if h.index < len(h.entries)-1 {
		h.entries = h.entries[:h.index+1]
	}
	h.entries = append(h.entries, path)
	h.index = len(h.entries) - 1
}

// Back moves back and returns path if possible.
func (h *History) Back() (string, bool) {
	if h == nil || h.index <= 0 {
		return "", false
	}
	h.index--
	return h.entries[h.index], true
}

// Forward moves forward and returns path if possible.
func (h *History) Forward() (string, bool) {
	if h == nil || h.index >= len(h.entries)-1 {
		return "", false
	}
	h.index++
	return h.entries[h.index], true
}

// Entries returns all history entries.
func (h *History) Entries() []string {
	return append([]string(nil), h.entries...)
}
