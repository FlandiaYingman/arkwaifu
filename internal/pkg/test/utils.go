package test

import (
	"bytes"
	"crypto/sha256"
	"encoding/base32"
	"encoding/base64"
	"encoding/gob"
	"hash"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/flandiayingman/arkwaifu/internal/pkg/util/pathutil"
	"github.com/pkg/errors"
	"golang.org/x/mod/sumdb/dirhash"
)

func HashDir(dir string) string {
	hashDir, _ := dirhash.HashDir(dir, "", dirhash.DefaultHash)
	return hashDir
}

func HashFile(file string) string {
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		panic(err)
	}

	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func HashObj(obj interface{}) string {
	h := hashObj(obj)
	return base64.RawStdEncoding.EncodeToString(h.Sum(nil))
}

func HashObjSafe(obj interface{}) string {
	h := hashObj(obj)
	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(h.Sum(nil))
}

func hashObj(obj interface{}) hash.Hash {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(obj)
	if err != nil {
		panic(err)
	}

	h := sha256.New()
	_, err = io.Copy(h, &buf)
	if err != nil {
		panic(err)
	}
	return h
}

// AssertAllIn checks whether all files in dir A exists in dir B.
// This requires the file hierarchies, names and contents to be the same.
func AssertAllIn(dirA, dirB string) error {
	return filepath.WalkDir(dirA, func(srcPath string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}

		dstPath := pathutil.MustChangeParent(srcPath, dirA, dirB)

		srcHash := HashFile(srcPath)
		dstHash := HashFile(dstPath)

		if srcHash != dstHash {
			return errors.Errorf("src %v hash %v != dst %v hash %v", srcPath, srcHash, dstPath, dstHash)
		}
		return nil
	})
}
