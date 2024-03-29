package arkassets

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
)

// Version represents a raw response of "https://ak-conf.hypergryph.com/config/prod/official/Android/version".
type Version struct {
	ResVersion    string `json:"resVersion"`
	ClientVersion string `json:"clientVersion"`
}

// HotUpdateList represents a raw response of "https://ak.hycdn.cn/assetbundle/official/Android/assets/{resVersion}/hot_update_list.json"
type HotUpdateList struct {
	FullPack        FullPack   `json:"fullPack"`
	VersionID       string     `json:"versionId"`
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
	MD5       string `json:"md5"`
	TotalSize int    `json:"totalSize"`
	AbSize    int    `json:"abSize"`
	Cid       int    `json:"cid"`
	Pid       string `json:"pid,omitempty"`
	Type      string `json:"type,omitempty"`
}

type PackInfo struct {
	Name      string `json:"name"`
	Hash      string `json:"hash"`
	MD5       string `json:"md5"`
	TotalSize int    `json:"totalSize"`
	AbSize    int    `json:"abSize"`
	Cid       int    `json:"cid"`
}

// GetRawVersion gets the raw response from Arknights version API.
func GetRawVersion() (Version, error) {
	resp, err := http.Get("https://ak-conf.hypergryph.com/config/prod/official/Android/version")
	if err != nil {
		return Version{}, errors.WithStack(err)
	}
	defer resp.Body.Close()

	arkVersion := Version{}
	err = json.NewDecoder(resp.Body).Decode(&arkVersion)
	if err != nil {
		return Version{}, errors.WithStack(err)
	}
	return arkVersion, nil
}

// GetRawResources gets the raw response of Arknights resource API with specified resource version.
func GetRawResources(resVersion string) (HotUpdateList, error) {
	urlResourceList := GetResURL(resVersion, "hot_update_list.json")

	resp, err := http.Get(urlResourceList)
	if err != nil {
		return HotUpdateList{}, errors.WithStack(err)
	}
	defer resp.Body.Close()

	resourcesList := HotUpdateList{}
	err = json.NewDecoder(resp.Body).Decode(&resourcesList)
	if err != nil {
		return HotUpdateList{}, errors.WithStack(err)
	}
	return resourcesList, nil
}

// GetResURL gets the URL to download the specified asset with specified resource version.
func GetResURL(resVersion string, res string) string {
	return fmt.Sprintf("https://ak.hycdn.cn/assetbundle/official/Android/assets/%v/%v", resVersion, res)
}
