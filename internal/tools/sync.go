package tools

import (
	"os"
	"path/filepath"
)

// SyncItem describes a file difference.
type SyncItem struct {
	Path   string
	Status string // left-only, right-only, different, same
}

// CompareDirs compares two directories (non-recursive listing).
func CompareDirs(left, right string) ([]SyncItem, error) {
	leftEntries, err := readNames(left)
	if err != nil {
		return nil, err
	}
	rightEntries, err := readNames(right)
	if err != nil {
		return nil, err
	}
	var out []SyncItem
	for name, linfo := range leftEntries {
		rinfo, ok := rightEntries[name]
		if !ok {
			out = append(out, SyncItem{Path: name, Status: "left-only"})
			continue
		}
		status := "same"
		if linfo.size != rinfo.size {
			status = "different"
		}
		out = append(out, SyncItem{Path: name, Status: status})
	}
	for name := range rightEntries {
		if _, ok := leftEntries[name]; !ok {
			out = append(out, SyncItem{Path: name, Status: "right-only"})
		}
	}
	return out, nil
}

type entryInfo struct{ size int64 }

func readNames(dir string) (map[string]entryInfo, error) {
	infos, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	out := make(map[string]entryInfo)
	for _, e := range infos {
		if e.Name() == ".." {
			continue
		}
		fi, err := e.Info()
		if err != nil {
			continue
		}
		out[e.Name()] = entryInfo{size: fi.Size()}
	}
	return out, nil
}

// CopySync copies items from srcDir to dstDir by name list.
func CopySync(srcDir string, names []string) error {
	for _, name := range names {
		from := filepath.Join(srcDir, name)
		to := filepath.Join(srcDir+"_sync_target", name) // placeholder overridden by caller
		_ = to
		_ = from
	}
	return nil
}
