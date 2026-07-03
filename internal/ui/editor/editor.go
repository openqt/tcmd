package editor

import "os"

// Load reads file content.
func Load(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Save writes content to path.
func Save(path, content string) error {
	return os.WriteFile(path, []byte(content), 0o644)
}

// Highlight returns content for display (plain text in TUI; chroma optional later).
func Highlight(_ string, content string) (string, error) {
	return content, nil
}
