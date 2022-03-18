// Package arkres provides functions to retrieve resources of the game Arknights.
// This package only do fetching, unzipping, unpacking and decrypting resources.
// To do other works, such as parsing BSON or JSON tables or processing graphics assets, use the subpackages.
//
// tools/extractor is required.
package arkres

import (
	"context"
	"github.com/ahmetb/go-linq/v3"
	"github.com/pkg/errors"
	"os"
	"regexp"
)

const (
	DefaultPrefix = "assets/torappu/dynamicassets"
)

func Get(ctx context.Context, resVersion string, dst string, filter *regexp.Regexp) error {
	infos, err := GetResInfos(resVersion)
	if err != nil {
		return errors.WithStack(err)
	}

	infos = FilterResInfosRegexp(infos, filter)

	err = os.RemoveAll(dst)
	if err != nil {
		return errors.WithStack(err)
	}
	err = GetRes(ctx, infos, dst)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func Update(ctx context.Context, oldResVer string, newResVer string, dst string, filter *regexp.Regexp) error {
	oldResInfos, err := GetResInfos(oldResVer)
	if err != nil {
		return errors.WithStack(err)
	}
	newResInfos, err := GetResInfos(newResVer)
	if err != nil {
		return errors.WithStack(err)
	}

	oldResInfos = FilterResInfosRegexp(oldResInfos, filter)
	newResInfos = FilterResInfosRegexp(newResInfos, filter)

	infos := calcDiff(newResInfos, oldResInfos)
	err = GetRes(ctx, infos, dst)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func SetChatMask(newChatMask []byte) {
	setChatMask(newChatMask)
}

func GetLatestResVersion() (string, error) {
	return GetResVersion()
}

// calcDiff returns the infos in "new" but not in "old".
// It checks only "Name" and "MD5" for equality.
// The returned res.Info have the same "resVersion" as the res.Info in "new".
func calcDiff(new []Info, old []Info) []Info {
	convert := func(i interface{}) interface{} {
		info := i.(Info)
		return struct {
			name string
			md5  string
		}{
			name: info.Name,
			md5:  info.MD5,
		}
	}
	newQuery := linq.From(new)
	oldQuery := linq.From(old)
	var result []Info
	newQuery.
		ExceptBy(oldQuery, convert).
		ToSlice(&result)
	return result
}
