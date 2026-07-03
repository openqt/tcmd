package filesrc

import (
	"io"
	"os"
	"time"
)

// Entry describes one directory listing item.
type Entry struct {
	Name    string
	Path    string
	IsDir   bool
	Size    int64
	ModTime time.Time
	Mode    os.FileMode
}

// Source reads directory listings from a backing store.
type Source interface {
	List(path string, showHidden bool) ([]Entry, error)
	Stat(path string) (os.FileInfo, error)
	OpenRead(path string) (io.ReadCloser, error)
	Mkdir(path string) error
	Remove(path string) error
	RemoveAll(path string) error
	Rename(oldPath, newPath string) error
	CopyFile(src, dst string) error
}
