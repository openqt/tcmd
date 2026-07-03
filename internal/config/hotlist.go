package config

import (
	"os"
	"path/filepath"
	"strings"
)

// Hotlist stores favorite directories.
type Hotlist struct {
	Paths []string
}

// LoadHotlist reads dirhotlist.txt from config dir.
func LoadHotlist() (*Hotlist, error) {
	dir, err := EnsureDir()
	if err != nil {
		return nil, err
	}
	path := filepath.Join(dir, "dirhotlist.txt")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Hotlist{}, nil
		}
		return nil, err
	}
	lines := strings.Split(string(data), "\n")
	h := &Hotlist{}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			h.Paths = append(h.Paths, line)
		}
	}
	return h, nil
}

// SaveHotlist writes dirhotlist.txt.
func (h *Hotlist) Save() error {
	dir, err := EnsureDir()
	if err != nil {
		return err
	}
	path := filepath.Join(dir, "dirhotlist.txt")
	var b strings.Builder
	for _, p := range h.Paths {
		b.WriteString(p)
		b.WriteString("\n")
	}
	return os.WriteFile(path, []byte(b.String()), 0o644)
}

// Add appends a unique path.
func (h *Hotlist) Add(path string) {
	for _, p := range h.Paths {
		if p == path {
			return
		}
	}
	h.Paths = append(h.Paths, path)
}
