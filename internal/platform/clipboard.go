package platform

import (
	"github.com/atotto/clipboard"
)

// SetClipboard writes text to the system clipboard.
func SetClipboard(text string) error {
	return clipboard.WriteAll(text)
}

// GetClipboard reads text from the system clipboard.
func GetClipboard() (string, error) {
	return clipboard.ReadAll()
}
