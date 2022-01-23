package arkres

import (
	"fmt"
	"golang.org/x/mod/sumdb/dirhash"
	"os"
	"path/filepath"
	"testing"
)

func hashGetArt(t *testing.T, resVersion string) string {
	dst := filepath.Clean(fmt.Sprintf("tmp/ARKRES_TEST-%v", resVersion))
	_ = os.RemoveAll(dst)
	_ = os.MkdirAll(dst, os.ModePerm)

	err := GetArt(dst, resVersion)
	if err != nil {
		t.Fatalf("err %v", err)
	}

	hash, err := dirhash.HashDir(dst, "", dirhash.DefaultHash)
	if err != nil {
		t.Fatalf("err %v", err)
	}

	_ = os.RemoveAll(dst)
	return hash
}

func TestGetArt(t *testing.T) {
	tests := []struct {
		name       string
		resVersion string
		want       string
	}{
		{
			name:       "[CN UPDATE] Client:1.7.01 Data:22-01-14-17-58-34-bd68ad",
			resVersion: "22-01-14-17-58-34-bd68ad",
			want:       "h1:R6xBu0bIzs4Kx1hKAO63f6bku9PTdoX4XQAs6K99ZZo=",
		},
		{
			name:       "[CN UPDATE] Client:1.6.01 Data:21-12-01-03-53-27-2e01ea",
			resVersion: "21-12-01-03-53-27-2e01ea",
			want:       "h1:sj8/tcwxidtHiFEr3lJLiq9ztwwvLQmrFnIqXRk9WSo=",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := hashGetArt(t, test.resVersion)
			if got != test.want {
				t.Fatalf("got %v, but want %v", got, test.want)
			}
		})
	}
}
