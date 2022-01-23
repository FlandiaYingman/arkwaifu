package arkres

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	versionUrl = "https://ak-conf.hypergryph.com/config/prod/official/Android/version"
)

func getVersionUrl() string {
	return versionUrl
}

func getAssetBundleUrl(resVersion string, assetBundle string) string {
	return fmt.Sprintf("https://ak.hycdn.cn/assetbundle/official/Android/assets/%v/%v", resVersion, assetBundle)
}

func getCurrentArkVersion() (Version, error) {
	resp, err := http.Get(getVersionUrl())
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

func getArkResources(resVersion string) (Resources, error) {
	urlResourceList := getAssetBundleUrl(resVersion, "hot_update_list.json")

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

func getCurrentArkResources() (Version, Resources, error) {
	version, err := getCurrentArkVersion()
	if err != nil {
		return Version{}, Resources{}, err
	}
	resources, err := getArkResources(version.ResVersion)
	if err != nil {
		return Version{}, Resources{}, err
	}
	return version, resources, nil
}
