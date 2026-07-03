package tools

import (
	"os"
	"path/filepath"
	"strings"
)

// FindOptions configures a file search.
type FindOptions struct {
	Root     string
	NameMask string
	Content  string
}

// Find walks root and returns matching file paths.
func Find(opt FindOptions) ([]string, error) {
	var out []string
	err := filepath.Walk(opt.Root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		base := info.Name()
		if opt.NameMask != "" && !matchSimple(base, opt.NameMask) {
			return nil
		}
		if opt.Content != "" {
			data, err := os.ReadFile(path)
			if err != nil || !strings.Contains(string(data), opt.Content) {
				return nil
			}
		}
		out = append(out, path)
		return nil
	})
	return out, err
}

func matchSimple(name, mask string) bool {
	mask = strings.ToLower(mask)
	name = strings.ToLower(name)
	if strings.Contains(mask, "*") {
		parts := strings.Split(mask, "*")
		if len(parts) == 2 {
			return strings.HasPrefix(name, parts[0]) && strings.HasSuffix(name, parts[1])
		}
	}
	return strings.Contains(name, mask)
}
