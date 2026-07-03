package differ

import (
	"os"
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
)

// Compare returns a unified diff style text for two files.
func Compare(leftPath, rightPath string) (string, error) {
	l, err := os.ReadFile(leftPath)
	if err != nil {
		return "", err
	}
	r, err := os.ReadFile(rightPath)
	if err != nil {
		return "", err
	}
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(string(l), string(r), false)
	return dmp.DiffPrettyText(diffs), nil
}

// CompareText diffs two strings.
func CompareText(a, b string) string {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(a, b, false)
	text := dmp.DiffPrettyText(diffs)
	return strings.TrimSpace(text)
}
