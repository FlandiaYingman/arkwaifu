package arkassets

import (
	"context"
	"fmt"
	"github.com/ahmetb/go-linq/v3"
	"github.com/flandiayingman/arkwaifu/internal/pkg/ark"
	"github.com/flandiayingman/arkwaifu/internal/pkg/cols"
	"github.com/flandiayingman/arkwaifu/internal/pkg/util/fileutil"
	"github.com/flandiayingman/arkwaifu/internal/pkg/util/pathutil"
	"github.com/pkg/errors"
	"regexp"
	"strings"
)

const (
	AssetBaseUrl = "https://ak.hycdn.cn/assetbundle/official/Android/assets"
)

type Info struct {
	Name    string
	MD5     string
	Version ark.Version
}

func (i Info) Url() string {
	i.Name = strings.ReplaceAll(i.Name, `/`, `_`)
	i.Name = strings.ReplaceAll(i.Name, `#`, `__`)
	i.Name = regexp.MustCompile(`\.(.*?)$`).ReplaceAllString(i.Name, ".dat")
	return fmt.Sprintf("%s/%s/%s", AssetBaseUrl, i.Version, i.Name)
}

func GetGameAssets(ctx context.Context, version ark.Version, dst string, patterns []string) error {
	var err error

	if version == "" {
		version, err = GetLatestVersion()
		if err != nil {
			return errors.WithStack(err)
		}
	}

	infoList, err := GetInfoList(version)
	if err != nil {
		return errors.WithStack(err)
	}

	infoList, err = filter(infoList, patterns)
	if err != nil {
		return errors.WithStack(err)
	}

	err = get(ctx, infoList, dst)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
func UpdateGameAssets(ctx context.Context, oldResVer string, newResVer string, dst string, patterns []string) error {
	var err error

	if oldResVer == "" {
		return GetGameAssets(ctx, newResVer, dst, patterns)
	}
	if newResVer == "" {
		newResVer, err = GetLatestVersion()
		if err != nil {
			return errors.WithStack(err)
		}
	}

	oldInfoList, err := GetInfoList(oldResVer)
	if err != nil {
		return errors.WithStack(err)
	}
	newInfoList, err := GetInfoList(newResVer)
	if err != nil {
		return errors.WithStack(err)
	}

	oldInfoList, err = filter(oldInfoList, patterns)
	if err != nil {
		return errors.WithStack(err)
	}
	newInfoList, err = filter(newInfoList, patterns)
	if err != nil {
		return errors.WithStack(err)
	}

	infoList := CalculateDifferences(newInfoList, oldInfoList)
	err = get(ctx, infoList, dst)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func GetLatestVersion() (string, error) {
	version, err := GetRawVersion()
	if err != nil {
		return "", errors.WithStack(err)
	}

	return version.ResVersion, nil
}

func GetInfoList(version ark.Version) ([]Info, error) {
	resources, err := GetRawResources(version)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	infos := make([]Info, len(resources.AbInfos))
	for i, info := range resources.AbInfos {
		if err != nil {
			return nil, errors.WithStack(err)
		}
		infos[i] = Info{
			Name:    info.Name,
			MD5:     info.MD5,
			Version: version,
		}
	}
	return infos, nil
}

func filter(infoList []Info, patterns []string) ([]Info, error) {
	return cols.FilterErr(infoList, func(i Info) (bool, error) {
		return pathutil.MatchAny(patterns, i.Name)
	})
}

func get(ctx context.Context, infoList []Info, dst string) error {
	fetch, err := fetch(ctx, infoList)
	if err != nil {
		return errors.WithStack(err)
	}
	unzip, err := unzip(ctx, fetch)
	if err != nil {
		return errors.WithStack(err)
	}
	unpack, err := unpack(ctx, unzip)
	if err != nil {
		return errors.WithStack(err)
	}
	err = fileutil.MoveAllContent(unpack, dst)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func CalculateDifferences(new []Info, old []Info) []Info {
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
