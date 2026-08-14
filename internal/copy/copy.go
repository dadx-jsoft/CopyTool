package copy

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type FilterMode int

const (
	FilterAll FilterMode = iota
	FilterInclude
	FilterExclude
)

// CopyFilesToDirectory copies files from sourceDir into destDir flattened.
// FilterAll copies everything. FilterInclude/FilterExclude use exts
// (case-insensitive, with or without leading dot).
func CopyFilesToDirectory(sourceDir, destDir string, mode FilterMode, exts []string) (int64, error) {
	extSet := makeExtSet(exts)

	entries, err := os.ReadDir(sourceDir)
	if err != nil {
		return -1, err
	}

	var count int64
	for _, entry := range entries {
		srcPath := filepath.Join(sourceDir, entry.Name())
		if entry.IsDir() {
			n, err := CopyFilesToDirectory(srcPath, destDir, mode, exts)
			if err != nil {
				return -1, err
			}
			count += n
			continue
		}

		if !shouldCopy(entry.Name(), mode, extSet) {
			continue
		}

		destPath := filepath.Join(destDir, entry.Name())
		if _, err := os.Stat(destPath); err == nil {
			destPath = filepath.Join(destDir, fmt.Sprintf("%d-%s", time.Now().UnixMilli(), entry.Name()))
		} else if err != nil && !os.IsNotExist(err) {
			return -1, err
		}

		if err := copyFile(srcPath, destPath); err != nil {
			return -1, err
		}
		count++
	}
	return count, nil
}

func makeExtSet(exts []string) map[string]struct{} {
	if len(exts) == 0 {
		return nil
	}
	set := make(map[string]struct{}, len(exts))
	for _, e := range exts {
		e = normalizeExt(e)
		if e != "" {
			set[e] = struct{}{}
		}
	}
	if len(set) == 0 {
		return nil
	}
	return set
}

func normalizeExt(ext string) string {
	ext = strings.TrimSpace(strings.ToLower(ext))
	ext = strings.TrimPrefix(ext, ".")
	return ext
}

func ParseExtensions(input string) []string {
	parts := strings.FieldsFunc(input, func(r rune) bool {
		return r == ',' || r == ';' || r == ' ' || r == '\t'
	})
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if e := normalizeExt(p); e != "" {
			out = append(out, e)
		}
	}
	return out
}

func fileExt(name string) string {
	return strings.TrimPrefix(strings.ToLower(filepath.Ext(name)), ".")
}

func shouldCopy(name string, mode FilterMode, extSet map[string]struct{}) bool {
	switch mode {
	case FilterInclude:
		if extSet == nil {
			return false
		}
		_, ok := extSet[fileExt(name)]
		return ok
	case FilterExclude:
		if extSet == nil {
			return true
		}
		_, excluded := extSet[fileExt(name)]
		return !excluded
	default:
		return true
	}
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}

	_, copyErr := io.Copy(out, in)
	closeErr := out.Close()
	if copyErr != nil {
		return copyErr
	}
	return closeErr
}
