package operations

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/openqt/tcmd/internal/filesrc"
)

// CopyEntries copies files and directories into targetDir.
func CopyEntries(src filesrc.Source, entries []filesrc.Entry, targetDir string) error {
	for _, e := range entries {
		dst := filepath.Join(targetDir, filepath.Base(e.Path))
		if e.IsDir {
			if err := copyDir(src, e.Path, dst); err != nil {
				return err
			}
			continue
		}
		if err := copyFile(src, e.Path, dst); err != nil {
			return err
		}
	}
	return nil
}

// MoveEntries renames entries into targetDir.
func MoveEntries(src filesrc.Source, entries []filesrc.Entry, targetDir string) error {
	for _, e := range entries {
		dst := filepath.Join(targetDir, filepath.Base(e.Path))
		if err := src.Rename(e.Path, dst); err != nil {
			return err
		}
	}
	return nil
}

// DeleteEntries removes entries permanently.
func DeleteEntries(src filesrc.Source, entries []filesrc.Entry) error {
	for _, e := range entries {
		if e.IsDir {
			if err := src.RemoveAll(e.Path); err != nil {
				return err
			}
			continue
		}
		if err := src.Remove(e.Path); err != nil {
			return err
		}
	}
	return nil
}

// Mkdir creates a directory under base.
func Mkdir(src filesrc.Source, base, name string) (string, error) {
	path := filepath.Join(base, name)
	if err := src.Mkdir(path); err != nil {
		return "", err
	}
	return path, nil
}

// Rename renames an entry.
func Rename(src filesrc.Source, entry filesrc.Entry, newName string) error {
	dst := filepath.Join(filepath.Dir(entry.Path), newName)
	return src.Rename(entry.Path, dst)
}

func copyFile(src filesrc.Source, from, to string) error {
	in, err := src.OpenRead(from)
	if err != nil {
		return err
	}
	defer in.Close()

	if err := os.MkdirAll(filepath.Dir(to), 0o755); err != nil {
		return err
	}
	out, err := os.Create(to)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Close()
}

func copyDir(src filesrc.Source, from, to string) error {
	if err := os.MkdirAll(to, 0o755); err != nil {
		return err
	}
	entries, err := src.List(from, true)
	if err != nil {
		return err
	}
	for _, e := range entries {
		if e.Name == ".." {
			continue
		}
		dst := filepath.Join(to, e.Name)
		if e.IsDir {
			if err := copyDir(src, e.Path, dst); err != nil {
				return err
			}
			continue
		}
		if err := copyFile(src, e.Path, dst); err != nil {
			return err
		}
	}
	return nil
}

// DirSize calculates total bytes in a directory tree.
func DirSize(src filesrc.Source, path string) (int64, error) {
	entries, err := src.List(path, true)
	if err != nil {
		return 0, err
	}
	var total int64
	for _, e := range entries {
		if e.Name == ".." {
			continue
		}
		if e.IsDir {
			sub, err := DirSize(src, e.Path)
			if err != nil {
				return 0, err
			}
			total += sub
			continue
		}
		total += e.Size
	}
	return total, nil
}

// FormatDirSize returns human-readable size.
func FormatDirSize(n int64) string {
	return fmt.Sprintf("%d bytes", n)
}
