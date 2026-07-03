package filesrc

import (
	"strings"

	"github.com/openqt/tcmd/internal/filesrc/archive"
)

// IsArchiveFile reports zip extension.
func IsArchiveFile(path string) bool {
	if _, _, ok := archive.ParseArchivePath(path); ok {
		return true
	}
	return archive.IsZipPath(path)
}

// EnterZip returns virtual path for a zip file.
func EnterZip(zipPath string) string {
	return archive.JoinArchivePath(zipPath, "")
}

// ArchiveFileName returns base zip path from virtual path.
func ArchiveFileName(path string) string {
	if arch, _, ok := archive.ParseArchivePath(path); ok {
		return arch
	}
	return path
}

// IsVirtualArchivePath reports # syntax.
func IsVirtualArchivePath(path string) bool {
	return strings.Contains(path, "#")
}
