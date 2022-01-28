package arkres

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// ResInfo represents an abstraction of AbInfo and PackInfo.
type ResInfo struct {
	Name       string
	Hash       string
	ResVersion string
}

// GetUrl gets the URL for downloading this ResInfo.
func (i ResInfo) GetUrl() string {
	extensionRegexp := regexp.MustCompile(`\..*?$`)
	name := i.Name
	name = strings.ReplaceAll(name, `/`, `_`)
	name = strings.ReplaceAll(name, `#`, `__`)
	name = extensionRegexp.ReplaceAllString(name, ".dat")
	return getAssetUrl(i.ResVersion, name)
}

// Equals tests if two ResInfo are equal. Iff the two ResInfo have the same Name and Hash, they are equal.
func (i ResInfo) Equals(j ResInfo) bool {
	return i.Name == j.Name && i.Hash == j.Hash
}

// GetResVersion gets the current resource version.
func GetResVersion() (string, error) {
	version, err := GetRawVersion()
	if err != nil {
		return "", err
	}

	return version.ResVersion, nil
}

// GetResInfos gets the resource infos of specified resource version.
func GetResInfos(resVersion string) ([]ResInfo, error) {
	resources, err := GetRawResources(resVersion)
	if err != nil {
		return nil, err
	}

	infos := make([]ResInfo, len(resources.AbInfos))
	for i, info := range resources.AbInfos {
		infos[i] = ResInfo{
			Name:       info.Name,
			Hash:       info.Hash,
			ResVersion: resVersion,
		}
	}
	return infos, nil
}

// GetImageRes gets the resources to specified destination directory.
func GetImageRes(infos []ResInfo, dest string) error {
	urls := make([]string, len(infos))
	for i, info := range infos {
		urls[i] = info.GetUrl()
	}

	dstTmp := filepath.Join(dest, "tmp")
	var err error

	dstDl := filepath.Join(dstTmp, "dl")
	err = downloadResources(urls, dstDl)
	if err != nil {
		return err
	}

	dstUz := filepath.Join(dstTmp, "uz")
	err = unzipResources(dstDl, dstUz)
	if err != nil {
		return err
	}

	dstUp := filepath.Clean(dest)
	err = unpackImageResources(dstUz, dstUp)
	if err != nil {
		return err
	}

	_ = os.RemoveAll(dstTmp)
	return nil
}

// GetTextRes gets the resources to specified destination directory.
func GetTextRes(infos []ResInfo, dest string) error {
	urls := make([]string, len(infos))
	for i, info := range infos {
		urls[i] = info.GetUrl()
	}

	dstTmp := filepath.Join(dest, "tmp")
	var err error

	dstDl := filepath.Join(dstTmp, "dl")
	err = downloadResources(urls, dstDl)
	if err != nil {
		return err
	}

	dstUz := filepath.Join(dstTmp, "uz")
	err = unzipResources(dstDl, dstUz)
	if err != nil {
		return err
	}

	dstUp := filepath.Clean(dest)
	err = unpackTextResources(dstUz, dstUp)
	if err != nil {
		return err
	}

	//_ = os.RemoveAll(dstTmp)
	return nil
}

// FilterResInfos returns a slice containing all ResInfo matching the specified predicate.
func FilterResInfos(infos []ResInfo, predicate func(i ResInfo) bool) []ResInfo {
	result := make([]ResInfo, 0)
	for _, info := range infos {
		if predicate(info) {
			result = append(result, info)
		}
	}
	return result
}

// FilterResInfosRegexp returns a slice containing all ResInfo matching the specified regexp.
func FilterResInfosRegexp(infos []ResInfo, regexp *regexp.Regexp) []ResInfo {
	return FilterResInfos(infos, func(i ResInfo) bool {
		return regexp.MatchString(i.Name)
	})
}
