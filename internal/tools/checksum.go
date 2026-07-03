package tools

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

// ChecksumAlgo supported algorithms.
type ChecksumAlgo string

const (
	AlgoMD5    ChecksumAlgo = "md5"
	AlgoSHA1   ChecksumAlgo = "sha1"
	AlgoSHA256 ChecksumAlgo = "sha256"
)

// CalcChecksum hashes a file.
func CalcChecksum(path string, algo ChecksumAlgo) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	var h interface {
		io.Writer
		Sum([]byte) []byte
	}
	switch algo {
	case AlgoSHA1:
		h = sha1.New()
	case AlgoSHA256:
		h = sha256.New()
	default:
		h = md5.New()
	}
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// SplitFile splits file into parts with suffix .001, .002...
func SplitFile(path string, partSize int64) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if partSize <= 0 {
		partSize = 1024 * 1024
	}
	var parts []string
	for i := 0; i*int(partSize) < len(data); i++ {
		start := i * int(partSize)
		end := start + int(partSize)
		if end > len(data) {
			end = len(data)
		}
		partPath := fmt.Sprintf("%s.%03d", path, i+1)
		if err := os.WriteFile(partPath, data[start:end], 0o644); err != nil {
			return nil, err
		}
		parts = append(parts, partPath)
	}
	return parts, nil
}

// CombineFiles merges split parts back into target.
func CombineFiles(parts []string, target string) error {
	out, err := os.Create(target)
	if err != nil {
		return err
	}
	defer out.Close()
	for _, p := range parts {
		data, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		if _, err := out.Write(data); err != nil {
			return err
		}
	}
	return out.Close()
}
