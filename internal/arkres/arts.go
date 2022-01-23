package arkres

import (
	"github.com/thoas/go-funk"
	"os"
	"path/filepath"
	"strings"
)

var (
	ArtPrefixes = []string{
		"avg/bg",
		"avg/imgs",
	}
)

func GetLatestArt(dst string) error {
	version, err := getCurrentArkVersion()
	if err != nil {
		return err
	}

	return GetArt(dst, version.ResVersion)
}

func GetArt(dst string, resVersion string) error {
	resources, err := getArkResources(resVersion)

	resInfos := filterPrefix(resources.AbInfos, ArtPrefixes)
	resUrls := toResUrls(resInfos, resVersion)

	dstDownload := filepath.Join(dst, "download")
	err = downloadResources(resUrls, dstDownload)
	if err != nil {
		return err
	}

	dstUnzip := filepath.Join(dst, "unzip")
	err = unzipResources(dstDownload, dstUnzip)
	if err != nil {
		return err
	}

	dstUnpack := filepath.Clean(dst)
	err = unpackResources(dstUnzip, dstUnpack)
	if err != nil {
		return err
	}

	_ = os.RemoveAll(dstDownload)
	_ = os.RemoveAll(dstUnzip)
	return nil
}

func filterPrefix(abInfos []AbInfo, prefixes []string) []AbInfo {
	return funk.Filter(abInfos, func(i AbInfo) bool {
		for _, prefix := range prefixes {
			if strings.HasPrefix(i.Name, prefix) {
				return true
			}
		}
		return false
	}).([]AbInfo)
}

func toResUrls(infos []AbInfo, resVersion string) []string {
	return funk.Map(infos, func(i AbInfo) string {
		return i.toUrl(resVersion)
	}).([]string)
}
