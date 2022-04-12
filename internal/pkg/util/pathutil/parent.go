package pathutil

import (
	"fmt"
	"path/filepath"
)

// ChangeParent changes the parent of srcPath from srcDir to dstDir, and returns the dstPath.
// If srcPath is not a child of srcDir, it returns an error.
//
// For example
//	ChangeParent("/home/user/dir/file.txt", "/home", "/root") // -> "/root/user/dir/file.txt"
// 	ChangeParent("/home/user/file.txt", "/home/user", "/home/user/dir") // -> "/home/user/dir/file.txt"
func ChangeParent(srcPath string, srcDir string, dstDir string) (dstPath string, err error) {
	relativePath, err := filepath.Rel(srcDir, srcPath)
	if err != nil {
		err = fmt.Errorf("failed to make relative path: %w", err)
		return
	}
	dstPath = filepath.Join(dstDir, relativePath)
	return
}

// MustChangeParent calls ChangeParent and panics if an error occurs.
// It is safe to use MustChangeParent in, for example, os.WalkDir.
func MustChangeParent(srcPath string, srcDir string, dstDir string) string {
	dstPath, err := ChangeParent(srcPath, srcDir, dstDir)
	if err != nil {
		panic(err)
	}
	return dstPath
}
