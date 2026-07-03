package local

import (
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/openqt/tcmd/internal/filesrc"
)

// Source implements filesrc.Source for the local filesystem.
type Source struct{}

// New returns a local filesystem source.
func New() *Source {
	return &Source{}
}

// List returns directory entries with optional hidden files.
func (s *Source) List(path string, showHidden bool) ([]filesrc.Entry, error) {
	path = filepath.Clean(path)
	infos, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	entries := make([]filesrc.Entry, 0, len(infos)+1)
	parent := filepath.Dir(path)
	if parent != path {
		entries = append(entries, filesrc.Entry{
			Name:  "..",
			Path:  parent,
			IsDir: true,
		})
	}

	for _, info := range infos {
		name := info.Name()
		if !showHidden && strings.HasPrefix(name, ".") {
			continue
		}
		full := filepath.Join(path, name)
		detail, err := info.Info()
		if err != nil {
			continue
		}
		entries = append(entries, filesrc.Entry{
			Name:    name,
			Path:    full,
			IsDir:   detail.IsDir(),
			Size:    detail.Size(),
			ModTime: detail.ModTime(),
			Mode:    detail.Mode(),
		})
	}
	sort.SliceStable(entries[minParent(entries):], func(i, j int) bool {
		a, b := entries[i+minParent(entries)], entries[j+minParent(entries)]
		if a.IsDir != b.IsDir {
			return a.IsDir
		}
		return strings.ToLower(a.Name) < strings.ToLower(b.Name)
	})
	return entries, nil
}

func minParent(entries []filesrc.Entry) int {
	if len(entries) > 0 && entries[0].Name == ".." {
		return 1
	}
	return 0
}

// Stat returns file info.
func (s *Source) Stat(path string) (os.FileInfo, error) {
	return os.Stat(path)
}

// OpenRead opens a file for reading.
func (s *Source) OpenRead(path string) (io.ReadCloser, error) {
	return os.Open(path)
}

// Mkdir creates a directory.
func (s *Source) Mkdir(path string) error {
	return os.Mkdir(path, 0o755)
}

// Remove deletes a file.
func (s *Source) Remove(path string) error {
	return os.Remove(path)
}

// RemoveAll deletes a directory tree.
func (s *Source) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

// Rename moves or renames a path.
func (s *Source) Rename(oldPath, newPath string) error {
	return os.Rename(oldPath, newPath)
}

// CopyFile copies a single file.
func (s *Source) CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Close()
}
