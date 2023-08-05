package pathutil

import (
	"github.com/gobwas/glob"
	"io/fs"
	"path/filepath"
)

func MatchAny(patterns []string, name string) (bool, error) {
	name = filepath.ToSlash(name)
	for _, pattern := range patterns {
		match, err := Match(pattern, name)
		if err != nil {
			return false, err
		}
		if match {
			return true, nil
		}
	}
	return false, nil
}

func Match(pattern string, name string) (bool, error) {
	globber, err := glob.Compile(pattern, '/', '\\')
	if err != nil {
		return false, err
	}
	return globber.Match(name), nil
}

func Glob(pattern []string, root string) ([]string, error) {
	results := make([]string, 0)
	return results, filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		rel, err := filepath.Rel(root, p)
		if err != nil {
			return err
		}
		match, err := MatchAny(pattern, rel)
		if err != nil {
			return err
		}
		if match {
			results = append(results, p)
		}

		return nil
	})
}
