package arkres

import (
	"regexp"
	"strings"
)

type Version struct {
	ResVersion    string `json:"resVersion"`
	ClientVersion string `json:"clientVersion"`
}

type Resources struct {
	FullPack        FullPack   `json:"fullPack"`
	VersionId       string     `json:"versionId"`
	AbInfos         []AbInfo   `json:"abInfos"`
	CountOfTypedRes int        `json:"countOfTypedRes"`
	PackInfos       []PackInfo `json:"packInfos"`
}

type FullPack struct {
	TotalSize int    `json:"totalSize"`
	AbSize    int    `json:"abSize"`
	Type      string `json:"type"`
	Cid       int    `json:"cid"`
}

type AbInfo struct {
	Name      string `json:"name"`
	Hash      string `json:"hash"`
	Md5       string `json:"md5"`
	TotalSize int    `json:"totalSize"`
	AbSize    int    `json:"abSize"`
	Cid       int    `json:"cid"`
	Pid       string `json:"pid,omitempty"`
	Type      string `json:"type,omitempty"`
}

type PackInfo struct {
	Name      string `json:"name"`
	Hash      string `json:"hash"`
	Md5       string `json:"md5"`
	TotalSize int    `json:"totalSize"`
	AbSize    int    `json:"abSize"`
	Cid       int    `json:"cid"`
}

type ResInfo interface {
	toUrl(resVersion string) string
}

var extensionRegexp = regexp.MustCompile(`\..*?$`)

func (i AbInfo) toUrl(resVersion string) string {
	name := i.Name
	name = strings.ReplaceAll(name, `/`, `_`)
	name = strings.ReplaceAll(name, `#`, `__`)
	name = extensionRegexp.ReplaceAllString(name, ".dat")
	return getAssetBundleUrl(resVersion, name)
}

func (i PackInfo) toUrl(resVersion string) string {
	name := i.Name
	name = strings.ReplaceAll(name, `/`, `_`)
	name = strings.ReplaceAll(name, `#`, `__`)
	name = extensionRegexp.ReplaceAllString(name, ".dat")
	return getAssetBundleUrl(resVersion, name)
}
