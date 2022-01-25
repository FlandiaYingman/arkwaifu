package test

import (
	"bytes"
	"crypto/sha256"
	"encoding/base32"
	"encoding/base64"
	"encoding/gob"
	"golang.org/x/mod/sumdb/dirhash"
	"hash"
	"io"
)

func HashDir(dir string) string {
	hashDir, _ := dirhash.HashDir(dir, "", dirhash.DefaultHash)
	return hashDir
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
