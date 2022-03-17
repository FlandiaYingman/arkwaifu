package pathutil

import (
	"fmt"
	"path/filepath"
	"strings"
)

func ReplaceExt(path string, ext string) string {
	withoutExt := strings.TrimSuffix(path, filepath.Ext(path))
	return fmt.Sprintf("%v%v", withoutExt, ext)
}

func ReplaceParent(srcPath string, srcDir string, dstDir string) (dstPath string) {
	relativePath, err := filepath.Rel(srcDir, srcPath)
	if err != nil {
		return "."
	}
	return filepath.Join(dstDir, relativePath)
}

func ReplaceParentExt(srcPath, srcDir, dstDir, dstExt string) (dstPath string) {
	dstPath = ReplaceExt(ReplaceParent(srcPath, srcDir, dstDir), dstExt)
	return
}
