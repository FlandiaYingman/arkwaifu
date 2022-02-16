package arkres

import (
	"arkwaifu/internal/app/test"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"testing"
)

func TestGetResVersion(t *testing.T) {
	got, err := GetResVersion()
	if err != nil {
		t.Fatalf("GetResVersion() error = %v", err)
	}
	regexpResVersion := regexp.MustCompile(`(\d{2}-){6}\w{6}`)
	if !regexpResVersion.MatchString(got) {
		t.Errorf("GetResVersion() got = %v, regexp %v", got, regexpResVersion)
	}
}

func TestGetResInfos(t *testing.T) {
	type args struct {
		resVersion string
	}
	tests := []struct {
		name     string
		args     args
		wantHash string
		wantErr  bool
	}{
		{
			name:     "[CN UPDATE] Client:1.7.01 Data:22-01-14-17-58-34-bd68ad",
			args:     struct{ resVersion string }{"22-01-14-17-58-34-bd68ad"},
			wantHash: "4gTd0xUQIth2Ia5rl0tQiAMySpsxRKhb6ll1GPsdT6U",
			wantErr:  false,
		},
		{
			name:     "[CN UPDATE] Client:1.6.01 Data:21-12-01-03-53-27-2e01ea",
			args:     struct{ resVersion string }{"21-12-01-03-53-27-2e01ea"},
			wantHash: "C1Nu+RA6363c3GAb0tWFKdanh7A4P1xRSZbbUwLIwK0",
			wantErr:  false,
		},
		{
			name:     "404 URL",
			args:     struct{ resVersion string }{"00-00-00-00-00-00-000000"},
			wantHash: "",
			wantErr:  true,
		},
		{
			name:     "Invalid URL",
			args:     struct{ resVersion string }{`!@#$%^&*()-=_+/\/\`},
			wantHash: "",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetResInfos(tt.args.resVersion)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetResInfos() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if err != nil {
				return
			}
			gotHash := test.HashObj(got)
			if !reflect.DeepEqual(gotHash, tt.wantHash) {
				t.Errorf("GetResInfos() gotHash = %v, wantHash %v", gotHash, tt.wantHash)
			}
		})
	}
}

func TestGetRes(t *testing.T) {
	tmp := "tmp"
	_ = os.RemoveAll(tmp)
	tests := []struct {
		name       string
		resVersion string
		filter     *regexp.Regexp
		wantHash   string
		wantErr    bool
	}{
		{
			name:       "[CN UPDATE] Client:1.6.01 Data:21-12-01-03-53-27-2e01ea; avg/",
			resVersion: "21-12-01-03-53-27-2e01ea",
			filter:     regexp.MustCompile("^avg/"),
			wantHash:   "h1:fpc225ht+N20EnZ0XRGeFGf/YesqbYjLozaEzP5bopY=",
			wantErr:    false,
		},
		{
			name:       "[CN UPDATE] Client:1.6.01 Data:21-12-01-03-53-27-2e01ea; battle/ or hotupdate/",
			resVersion: "21-12-01-03-53-27-2e01ea",
			filter:     regexp.MustCompile("^(battle)|(hotupdate)/.*"),
			wantHash:   "h1:6rswLEPmCp8hfuXvViyi620Agps8Hf3g8Hk64y5vRDw=",
			wantErr:    false,
		},
		{
			name:       "[CN UPDATE] Client:1.6.01 Data:21-12-01-03-53-27-2e01ea; avg/items/ or config/leveloperaconfig/",
			resVersion: "21-12-01-03-53-27-2e01ea",
			filter:     regexp.MustCompile("^(avg/items/)|(config/leveloperaconfig/).*"),
			wantHash:   "h1:+IwwrEjSRH/iS71hMqJh0XF3ADg1CP9sUg6b+1GQWzo=",
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			infos, err := GetResInfos(tt.resVersion)
			if err != nil {
				t.Errorf("GetResInfos() error = %v", err)
				return
			}
			infos = FilterResInfosRegexp(infos, tt.filter)
			dst := filepath.Join(tmp, fmt.Sprintf("arkres-test-%v", test.HashObjSafe(tt.name)[:8]))
			err = GetRes(infos, dst)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRes() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if err != nil {
				return
			}
			gotHash := test.HashDir(dst)
			if !reflect.DeepEqual(gotHash, tt.wantHash) {
				t.Errorf("GetRes() gotHash = %v, wantHash %v", gotHash, tt.wantHash)
			}
		})
	}
	_ = os.RemoveAll(tmp)
}
