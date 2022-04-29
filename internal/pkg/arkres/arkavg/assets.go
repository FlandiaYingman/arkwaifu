package arkavg

import (
	"path/filepath"
)

type Asset struct {
	Name string
	Kind Kind
}
type Kind string

const (
	KindImage      Kind = "images"
	KindBackground Kind = "backgrounds"
	KindCharacter  Kind = "characters"
)

func (a *Asset) FilePath(resDir string, prefix string) string {
	return filepath.Join(resDir, prefix, "avg", string(a.Kind), a.Name)
}
