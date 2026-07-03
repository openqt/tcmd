package viewer

import (
	"encoding/hex"
	"fmt"
	"os"
	"strings"
)

// Mode is the viewer display mode.
type Mode int

const (
	ModeText Mode = iota
	ModeHex
	ModeBinary
)

// Render reads a file and formats it for TUI display.
func Render(path string, mode Mode, maxBytes int) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	if len(data) > maxBytes {
		data = data[:maxBytes]
	}
	switch mode {
	case ModeHex:
		return fmt.Sprintf("HEX %s\n%s", path, hex.Dump(data)), nil
	case ModeBinary:
		return fmt.Sprintf("BINARY %s (%d bytes)", path, len(data)), nil
	default:
		return fmt.Sprintf("TEXT %s\n%s", path, string(data)), nil
	}
}

// Wrap applies simple line wrapping.
func Wrap(s string, width int) string {
	if width <= 0 {
		return s
	}
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		if len(line) > width {
			lines[i] = line[:width] + "..."
		}
	}
	return strings.Join(lines, "\n")
}
