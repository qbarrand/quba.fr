package fileutils

import (
	"encoding/hex"
	"fmt"
	"hash/fnv"
	"io"
	"path/filepath"
	"strings"
)

func HashReader(r io.Reader) (string, error) {
	hasher := fnv.New32()

	if _, err := io.Copy(hasher, r); err != nil {
		return "", fmt.Errorf("could copy bytes into the hasher: %v", err)
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func PathWithHash(path, hash string) string {
	ext := filepath.Ext(path)

	// 4 times faster than fmt.Sprintf
	var sb strings.Builder

	sb.WriteString(path[:len(path)-len(ext)])
	sb.WriteRune('.')
	sb.WriteString(hash)
	sb.WriteString(ext)

	return sb.String()
}
