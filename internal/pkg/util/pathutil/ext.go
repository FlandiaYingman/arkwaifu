package pathutil

import (
	"fmt"
	"path/filepath"
	"strings"
)

// RemoveExt returns the path without the extension.
func RemoveExt(path string) string {
	return strings.TrimSuffix(path, filepath.Ext(path))
}

// RemoveAllExt returns the path without all extensions.
func RemoveAllExt(path string) string {
	return strings.TrimSuffix(path, AllExt(path))
}

// AllExt returns all extensions (e.g., ".tar.gz" instead of ".gz") from the file name (including ".").
func AllExt(path string) string {
	_, after, found := strings.Cut(filepath.Base(path), ".")
	if found {
		return "." + after
	} else {
		return ""
	}
}

// ReplaceExt replaces the extension of the given path with the given extension.
func ReplaceExt(path string, ext string) string {
	withoutExt := RemoveExt(path)
	return fmt.Sprintf("%v%v", withoutExt, ext)
}

// ReplaceAllExt replaces all extensions of the given path with the given extension.
func ReplaceAllExt(path string, ext string) string {
	withoutExt := RemoveAllExt(path)
	return fmt.Sprintf("%v%v", withoutExt, ext)
}
