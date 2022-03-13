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
