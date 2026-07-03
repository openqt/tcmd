package archive

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ZipSource lists contents of a zip archive as a virtual directory.
type ZipSource struct {
	ArchivePath string
}

// NewZip returns a zip archive source.
func NewZip(path string) *ZipSource {
	return &ZipSource{ArchivePath: path}
}

// ZipEntry is one item inside an archive listing.
type ZipEntry struct {
	Name    string
	Path    string
	IsDir   bool
	Size    int64
	ModTime time.Time
}

// List returns entries at inner path (use "" for root).
func (z *ZipSource) List(inner string, showHidden bool) ([]ZipEntry, error) {
	r, err := zip.OpenReader(z.ArchivePath)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	inner = strings.TrimPrefix(filepath.ToSlash(inner), "/")
	prefix := inner
	if prefix != "" && !strings.HasSuffix(prefix, "/") {
		prefix += "/"
	}

	seen := map[string]bool{}
	var entries []ZipEntry
	if inner != "" {
		parent := filepath.Dir(inner)
		if parent == "." {
			parent = ""
		}
		entries = append(entries, ZipEntry{Name: "..", Path: parent, IsDir: true})
	}

	for _, f := range r.File {
		name := filepath.ToSlash(f.Name)
		if prefix != "" && !strings.HasPrefix(name, prefix) {
			continue
		}
		rest := strings.TrimPrefix(name, prefix)
		if rest == "" {
			continue
		}
		part := rest
		if i := strings.Index(rest, "/"); i >= 0 {
			part = rest[:i]
		}
		if part == "" || seen[part] {
			continue
		}
		seen[part] = true
		isDir := strings.Contains(rest, "/") || strings.HasSuffix(f.Name, "/")
		if !showHidden && strings.HasPrefix(part, ".") {
			continue
		}
		entries = append(entries, ZipEntry{
			Name:    part,
			Path:    filepath.Join(inner, part),
			IsDir:   isDir,
			Size:    int64(f.UncompressedSize64),
			ModTime: f.Modified,
		})
	}
	return entries, nil
}

// ExtractAll unpacks archive to target directory.
func ExtractAll(archivePath, targetDir string) error {
	r, err := zip.OpenReader(archivePath)
	if err != nil {
		return err
	}
	defer r.Close()
	for _, f := range r.File {
		path := filepath.Join(targetDir, f.Name)
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(path, 0o755); err != nil {
				return err
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return err
		}
		rc, err := f.Open()
		if err != nil {
			return err
		}
		out, err := os.Create(path)
		if err != nil {
			rc.Close()
			return err
		}
		_, err = io.Copy(out, rc)
		out.Close()
		rc.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

// PackZip creates a zip from files in srcDir (names only).
func PackZip(targetZip, srcDir string, names []string) error {
	out, err := os.Create(targetZip)
	if err != nil {
		return err
	}
	defer out.Close()
	w := zip.NewWriter(out)
	defer w.Close()
	for _, name := range names {
		path := filepath.Join(srcDir, name)
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		f, err := w.Create(name)
		if err != nil {
			return err
		}
		if _, err := f.Write(data); err != nil {
			return err
		}
	}
	return nil
}

// IsZipPath reports .zip extension.
func IsZipPath(path string) bool {
	return strings.HasSuffix(strings.ToLower(path), ".zip")
}

// ParseArchivePath splits panel path into archive and inner.
// Format: /path/file.zip#inner/dir
func ParseArchivePath(path string) (archive, inner string, ok bool) {
	i := strings.Index(path, "#")
	if i < 0 {
		return "", "", false
	}
	return path[:i], path[i+1:], true
}

// JoinArchivePath builds virtual archive path.
func JoinArchivePath(archive, inner string) string {
	if inner == "" {
		return archive + "#"
	}
	return archive + "#" + inner
}

// ModTime helper
func modTime(t time.Time) time.Time {
	if t.IsZero() {
		return time.Now()
	}
	return t
}
