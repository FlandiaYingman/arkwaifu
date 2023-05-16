// Package agdapi provides functionalities to work with the resources of the game from the GitHub repository - ArknightsGameData (https://github.com/Kengxxiao/ArknightsGameData).

package agdapi

import (
	"context"
	"github.com/ahmetb/go-linq/v3"
	"github.com/flandiayingman/arkwaifu/internal/pkg/arkres/arkconsts"
	"github.com/flandiayingman/arkwaifu/internal/pkg/util/fileutil"
	"os"
	"regexp"
)

const (
	DefaultPrefix = "assets/torappu/dynamicassets"
)

type ResourceInfo struct {
	Name string
}

type ResourceVersion struct {
	GameServer      arkconsts.GameServer
	ClientVersion   string
	ResourceVersion string
	CommitSHA       string
}

// GetRes gets the resources specified by infos into dst.
func GetRes(ctx context.Context, resourceVersion ResourceVersion, resourceInfoList []ResourceInfo, dst string) error {
	fetchTempDir, err := os.MkdirTemp("", "arkwaifu-arkres-agdapi-fetch-*")
	if err != nil {
		return err
	}
	unzipTempDir, err := os.MkdirTemp("", "arkwaifu-arkres-agdapi-unzip-*")
	if err != nil {
		return err
	}

	zipballPath, err := fetchZipball(ctx, resourceVersion, fetchTempDir)
	if err != nil {
		return err
	}

	err = unzipZipball(ctx, zipballPath, resourceInfoList, resourceVersion, unzipTempDir)
	if err != nil {
		return err
	}

	return fileutil.CopyAllContent(unzipTempDir, dst)
}

// FilterResInfos returns a slice containing all ResourceInfo matching the specified predicate.
func FilterResInfos(infos []ResourceInfo, predicate func(i ResourceInfo) bool) []ResourceInfo {
	result := make([]ResourceInfo, 0)
	for _, info := range infos {
		if predicate(info) {
			result = append(result, info)
		}
	}
	return result
}

// FilterResInfosRegexp returns a slice containing all ResourceInfo matching the specified regexp.
func FilterResInfosRegexp(infos []ResourceInfo, r []*regexp.Regexp) []ResourceInfo {
	return FilterResInfos(infos, func(i ResourceInfo) bool {
		return linq.From(r).
			AnyWith(func(j any) bool {
				return j.(*regexp.Regexp).MatchString(i.Name)
			})
	})
}
