package arkres

import (
	"crypto/sha1"
	"fmt"
	"golang.org/x/mod/sumdb/dirhash"
	"os"
	"path/filepath"
	"testing"
)

func hashGet(t *testing.T, targets []string, resVersion string) string {
	dst := filepath.Clean(fmt.Sprintf("tmp/ARKRES_TEST-%v-%v", resVersion, hashStrings(targets)))
	_ = os.RemoveAll(dst)
	_ = os.MkdirAll(dst, os.ModePerm)

	err := Get(targets, dst, resVersion)
	if err != nil {
		t.Fatalf("err %v", err)
	}

	hash := hashDir(t, dst)
	_ = os.RemoveAll(dst)
	return hash
}

func hashDir(t *testing.T, dir string) string {
	hash, err := dirhash.HashDir(dir, "", dirhash.DefaultHash)
	if err != nil {
		t.Fatalf("err %v", err)
	}
	return hash
}

func hashStrings(strings []string) string {
	h := sha1.New()
	h.Write([]byte(fmt.Sprintf("%v", strings)))

	return fmt.Sprintf("%x", h.Sum(nil))
}

func TestGet(t *testing.T) {
	tests := []struct {
		name       string
		targets    []string
		resVersion string
		resHash    string
	}{
		//{
		//	name:       "[CN UPDATE] Client:1.6.01 Data:21-12-01-03-53-27-2e01ea",
		//	targets:    []string{},
		//	resVersion: "21-12-01-03-53-27-2e01ea",
		//	resHash:    "h1:DrZj9nvv/qkxdsbMYIBJKPlV31cjvubZ+IcDoXuWdOc=",
		//},
		{
			name:       "[CN UPDATE] Client:1.6.01 Data:21-12-01-03-53-27-2e01ea",
			targets:    []string{"avg"},
			resVersion: "21-12-01-03-53-27-2e01ea",
			resHash:    "h1:fpc225ht+N20EnZ0XRGeFGf/YesqbYjLozaEzP5bopY=",
		},
		{
			name:       "[CN UPDATE] Client:1.6.01 Data:21-12-01-03-53-27-2e01ea",
			targets:    []string{"battle", "hotupdate"},
			resVersion: "21-12-01-03-53-27-2e01ea",
			resHash:    "h1:6rswLEPmCp8hfuXvViyi620Agps8Hf3g8Hk64y5vRDw=",
		},
		{
			name:       "[CN UPDATE] Client:1.6.01 Data:21-12-01-03-53-27-2e01ea",
			targets:    []string{"avg/items", "config/leveloperaconfig"},
			resVersion: "21-12-01-03-53-27-2e01ea",
			resHash:    "h1:+IwwrEjSRH/iS71hMqJh0XF3ADg1CP9sUg6b+1GQWzo=",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := hashGet(t, test.targets, test.resVersion)
			if got != test.resHash {
				t.Fatalf("got %v, but want %v", got, test.resHash)
			}
		})
	}
}
