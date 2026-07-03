package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// RenameRule applies a pattern to a filename.
type RenameRule struct {
	Prefix string
	Suffix string
	Find   string
	Repl   string
}

// ApplyRename generates a new name from rule.
func ApplyRename(name string, rule RenameRule) string {
	ext := filepath.Ext(name)
	base := strings.TrimSuffix(name, ext)
	if rule.Find != "" {
		base = strings.ReplaceAll(base, rule.Find, rule.Repl)
	}
	return rule.Prefix + base + rule.Suffix + ext
}

// PreviewRenames returns old->new mapping.
func PreviewRenames(names []string, rule RenameRule) map[string]string {
	out := make(map[string]string, len(names))
	for _, n := range names {
		out[n] = ApplyRename(n, rule)
	}
	return out
}

// ApplyRenames renames files in dir according to rule.
func ApplyRenames(dir string, names []string, rule RenameRule) error {
	for _, name := range names {
		newName := ApplyRename(name, rule)
		if newName == name {
			continue
		}
		if err := os.Rename(filepath.Join(dir, name), filepath.Join(dir, newName)); err != nil {
			return err
		}
	}
	return nil
}

// FormatRenameLog returns a human-readable preview.
func FormatRenameLog(preview map[string]string) string {
	var b strings.Builder
	for old, newName := range preview {
		if old == newName {
			continue
		}
		fmt.Fprintf(&b, "%s -> %s\n", old, newName)
	}
	return b.String()
}
