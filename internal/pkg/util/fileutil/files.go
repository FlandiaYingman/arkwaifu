package fileutil

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/flandiayingman/arkwaifu/internal/pkg/util/pathutil"
)

// MkParents creates all parents of a file with mode 0755 (before unmask).
// If the parents of the file already exist, MkParents does nothing.
func MkParents(filePath string) error {
	err := os.MkdirAll(filepath.Dir(filePath), 0755)
	if err != nil {
		return fmt.Errorf("failed to create parent directories of %s: %v", filePath, err)
	}
	return err
}

// MkFile creates a file and its parents with respectively mode 0666 and 0755 (both before unmask).
// If the parents of the file already exist, MkFile does nothing.
// If the file already exist, MkFile truncates it.
func MkFile(filePath string) (*os.File, error) {
	err := MkParents(filePath)
	if err != nil {
		return nil, err
	}
	file, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file %s: %v", filePath, err)
	}
	return file, nil
}

// CopyAllContent copies all files from srcDir to dstDir recursively.
// Only files' name and content are guaranteed to be the same.
// Only files are copied, (empty) directories are ignored.
// Any existing files will be overridden.
func CopyAllContent(srcDir string, dstDir string) error {
	return filepath.WalkDir(srcDir, func(srcPath string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}

		dstPath := pathutil.MustChangeParent(srcPath, srcDir, dstDir)
		err = CopyContent(srcPath, dstPath)
		if err != nil {
			return err
		}

		return nil
	})
}

// MoveAllContent moves all files from srcDir to dstDir recursively.
// Only files' name and content are guaranteed to be the same.
// Only files are moved, (empty) directories are ignored.
// Any existing files will be overridden.
//
// Use MoveAllContent to prevent os.Rename to panic "invalid cross-device link".
func MoveAllContent(srcDir string, dstDir string) error {
	return filepath.WalkDir(srcDir, func(srcPath string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}

		dstPath := pathutil.MustChangeParent(srcPath, srcDir, dstDir)
		err = MoveContent(srcPath, dstPath)
		if err != nil {
			return err
		}

		return nil
	})
}

// CopyContent copies the file src to dst.
// Only the file's name and content are guaranteed to be the same.
// The parents will be created if they don't exist.
// The existing file will be overridden.
func CopyContent(src string, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open src %s: %w", src, err)
	}
	defer srcFile.Close()

	dstFile, err := MkFile(dst)
	if err != nil {
		return fmt.Errorf("failed to make dst %s: %w", dst, err)
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("failed to copy to dst %s from src %s: %w", dst, src, err)
	}

	return nil
}

// MoveContent moves the file src to dst.
// Only the file's name and content are guaranteed to be the same.
// The parents will be created if they don't exist.
// The existing file will be overridden.
func MoveContent(src string, dst string) error {
	err := os.Rename(src, dst)
	if err != nil {
		err = CopyContent(src, dst)
		if err != nil {
			return err
		}

		err = os.Remove(src)
		if err != nil {
			return fmt.Errorf("failed to remove src %s: %w", src, err)
		}
	}

	return nil
}

// Exists checks if the file or directory exists.
// If Exists cannot determinate whether the file or directory exists (e.g., permission error), it returns an error.
func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// LowercaseAll renames the directory to lowercase recursively.
// If a file or directory is already lower-cased, LowercaseAll does nothing to it.
func LowercaseAll(dirPath string) error {
	return filepath.WalkDir(dirPath, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		return Lowercase(filepath.Join(dirPath, path))
	})
}

// Lowercase renames the file or directory to lowercase.
// If the file or directory is already lower-cased, it does nothing.
func Lowercase(path string) error {
	lowerPath := filepath.Join(filepath.Dir(path), strings.ToLower(filepath.Base(path)))
	if path != lowerPath {
		return os.Rename(path, lowerPath)
	} else {
		return nil
	}
}

// List lists all files and directories under the directory dirPath (no recursive).
func List(dirPath string) ([]string, error) {
	dirEntries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	var all []string
	for _, dirEntry := range dirEntries {
		all = append(all, filepath.Join(dirPath, dirEntry.Name()))
	}
	return all, err
}

// ListFiles lists all files under the directory dirPath (no recursive).
func ListFiles(dirPath string) ([]string, error) {
	dirEntries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	var all []string
	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() {
			continue
		}
		all = append(all, filepath.Join(dirPath, dirEntry.Name()))
	}
	return all, err
}

// ListDirs lists all directories under the directory dirPath (no recursive).
func ListDirs(dirPath string) ([]string, error) {
	dirEntries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	var all []string
	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() {
			all = append(all, filepath.Join(dirPath, dirEntry.Name()))
		}
	}
	return all, err
}

// ListAll lists all files and directories under the directory dirPath (recursive).
func ListAll(dirPath string) ([]string, error) {
	var all []string
	err := filepath.WalkDir(dirPath, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		all = append(all, filepath.Join(dirPath, path))
		return nil
	})
	return all, err
}

// ListAllFiles lists all files under the directory dirPath (recursive).
func ListAllFiles(dirPath string) ([]string, error) {
	var allFiles []string
	err := filepath.WalkDir(dirPath, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		allFiles = append(allFiles, filepath.Join(dirPath, path))
		return nil
	})
	return allFiles, err
}

// ListALlDirs lists all directories under the directory dirPath (recursive).
func ListALlDirs(dirPath string) ([]string, error) {
	var allFiles []string
	err := filepath.WalkDir(dirPath, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			allFiles = append(allFiles, filepath.Join(dirPath, path))
		}
		return nil
	})
	return allFiles, err
}

// MkFileFromReader calls MkFile and then writes the content of the reader to the file.
//
// Note that it's the caller's responsibility to close the reader.
func MkFileFromReader(filePath string, r io.Reader) error {
	f, err := MkFile(filePath)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	_, err = io.Copy(f, r)
	return nil
}

// MkFileFromBytes calls MkFile and then writes the content of the bytes to the file.
func MkFileFromBytes(filePath string, b []byte) error {
	return MkFileFromReader(filePath, bytes.NewReader(b))
}
