package fileutil

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func MoveAllFileContent(srcDir string, dstDir string) error {
	return filepath.WalkDir(srcDir, func(srcPath string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}

		dstPath, err := ChangeParent(srcPath, srcDir, dstDir)
		if err != nil {
			return err
		}
		err = os.MkdirAll(filepath.Dir(dstPath), 0755)
		if err != nil {
			return fmt.Errorf("couldn't make directories: %w", err)
		}
		err = MoveFileContent(srcPath, dstPath)
		if err != nil {
			return err
		}

		return nil
	})
}

func ChangeParent(srcPath string, srcDir string, dstDir string) (dstPath string, err error) {
	relativePath, err := filepath.Rel(srcDir, srcPath)
	if err != nil {
		err = fmt.Errorf("couldn't make relative path: %w", err)
		return
	}
	dstPath = filepath.Join(dstDir, relativePath)
	return
}

func CopyFileContent(src string, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("couldn't open src %s: %w", src, err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("couldn't open dst %s: %w", dst, err)
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("couldn't copy to dst %s from src %s: %w", dst, src, err)
	}

	return nil
}

func MoveFileContent(src string, dst string) error {
	err := CopyFileContent(src, dst)
	if err != nil {
		return err
	}

	// if copy was successful, remove src
	err = os.Remove(src)
	if err != nil {
		return fmt.Errorf("couldn't remove src %s: %w", src, err)
	}

	return nil
}
