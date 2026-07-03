package platform

import (
	"os"
	"path/filepath"
	"sort"
)

// Drives returns roots available on this system.
func Drives() []string {
	home, _ := os.UserHomeDir()
	drives := []string{"/", home, "/tmp"}
	seen := map[string]bool{}
	out := make([]string, 0, len(drives))
	for _, d := range drives {
		d = filepath.Clean(d)
		if d == "" || seen[d] {
			continue
		}
		if _, err := os.Stat(d); err != nil {
			continue
		}
		seen[d] = true
		out = append(out, d)
	}
	sort.Strings(out)
	return out
}
