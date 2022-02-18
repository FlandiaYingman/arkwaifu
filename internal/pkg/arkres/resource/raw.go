package resource

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Version represents a raw response of "https://ak-conf.hypergryph.com/config/prod/official/Android/version".
type Version struct {
	ResVersion    string `json:"resVersion"`
	ClientVersion string `json:"clientVersion"`
}

// Resources represents a raw response of ".../hot_update_list.json"
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

// GetRawVersion gets the raw response from Arknights version API.
func GetRawVersion() (Version, error) {
	resp, err := http.Get("https://ak-conf.hypergryph.com/config/prod/official/Android/version")
	if err != nil {
		return Version{}, err
	}

	//Ignoring error of Body.Close(). IDK if this is a right practice. It will probably be changed in the future.
	//TODO: Evaluate whether ignoring the error is OK.
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	arkVersion := Version{}
	err = json.NewDecoder(resp.Body).Decode(&arkVersion)
	if err != nil {
		return Version{}, err
	}
	return arkVersion, nil
}

// GetRawResources gets the raw response of Arknights resource API with specified resource version.
func GetRawResources(resVersion string) (Resources, error) {
	urlResourceList := getAssetUrl(resVersion, "hot_update_list.json")

	resp, err := http.Get(urlResourceList)
	if err != nil {
		return Resources{}, err
	}

	//Ignoring error of Body.Close(). IDK if this is a right practice. It will probably be changed in the future.
	//TODO: Evaluate whether ignoring the error is OK.
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	resourcesList := Resources{}
	err = json.NewDecoder(resp.Body).Decode(&resourcesList)
	if err != nil {
		return Resources{}, err
	}
	return resourcesList, nil
}

// getAssetUrl gets the URL to download the specified asset with specified resource version.
func getAssetUrl(resVersion string, asset string) string {
	return fmt.Sprintf("https://ak.hycdn.cn/assetbundle/official/Android/assets/%v/%v", resVersion, asset)
}
