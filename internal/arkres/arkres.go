package arkres

import (
	"os"
	"path/filepath"
	"strings"
)

func Get(targets []string, dest string, resVersion string) error {
	resources, err := getArkResources(resVersion)
	if err != nil {
		return err
	}

	resInfos := filterPrefix(resources.AbInfos, targets)
	resUrls := toResUrls(resInfos, resVersion)
	dstTmp := filepath.Join(dest, "tmp")

	dstDl := filepath.Join(dstTmp, "dl")
	err = downloadResources(resUrls, dstDl)
	if err != nil {
		return err
	}

	dstUz := filepath.Join(dstTmp, "uz")
	err = unzipResources(dstDl, dstUz)
	if err != nil {
		return err
	}

	dstUp := filepath.Clean(dest)
	err = unpackResources(dstUz, dstUp)
	if err != nil {
		return err
	}

	_ = os.RemoveAll(dstTmp)
	return nil
}

func GetLatest(targets []string, dest string) error {
	version, err := getCurrentArkVersion()
	if err != nil {
		return err
	}

	return Get(targets, dest, version.ResVersion)
}

func filterPrefix(abInfos []AbInfo, prefixes []string) []AbInfo {
	result := make([]AbInfo, 0)
	for _, info := range abInfos {
		hasPrefix := false
		for _, prefix := range prefixes {
			if strings.HasPrefix(info.Name, prefix) {
				hasPrefix = true
				break
			}
		}
		if hasPrefix || len(prefixes) == 0 {
			result = append(result, info)
		}
	}
	return result
}

func toResUrls(infos []AbInfo, resVersion string) []string {
	result := make([]string, len(infos))
	for i, info := range infos {
		result[i] = info.toUrl(resVersion)
	}
	return result
}

var (
	avgPrefixes = []string{"avg"}
)

func GetAvg(dest string, resVersion string) error {
	return Get(avgPrefixes, dest, resVersion)
}

func GetLatestAvg(dest string) error {
	return GetLatest(avgPrefixes, dest)
}
