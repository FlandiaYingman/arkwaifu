package hgapi

import (
	"context"
	"os"
	"regexp"
	"strings"

	"github.com/ahmetb/go-linq/v3"
	"github.com/flandiayingman/arkwaifu/internal/pkg/util/fileutil"
	"github.com/pkg/errors"
)

type Info struct {
	Name       string
	MD5        string
	ResVersion string
}

func (i *Info) CreateURL() string {
	name := i.Name
	name = strings.ReplaceAll(name, `/`, `_`)
	name = strings.ReplaceAll(name, `#`, `__`)
	name = regexp.MustCompile(`\.(.*?)$`).ReplaceAllString(name, ".dat")
	return GetResURL(i.ResVersion, name)
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
func GetResInfos(resVersion string) ([]Info, error) {
	resources, err := GetRawResources(resVersion)
	if err != nil {
		return nil, err
	}

	infos := make([]Info, len(resources.AbInfos))
	for i, info := range resources.AbInfos {
		if err != nil {
			return nil, errors.WithStack(err)
		}
		infos[i] = Info{
			Name:       info.Name,
			MD5:        info.MD5,
			ResVersion: resVersion,
		}
	}
	return infos, nil
}

// GetRes gets the resources specified by infos into dst.
func GetRes(ctx context.Context, infos []Info, dst string) error {
	err := preCheck()
	if err != nil {
		return err
	}

	tmpFetch, err := fetchRes(ctx, infos)
	if err != nil {
		return err
	}
	defer func() {
		_ = os.RemoveAll(tmpFetch)
	}()

	tmpUnzip, err := unzipRes(ctx, tmpFetch)
	if err != nil {
		return err
	}
	defer func() {
		_ = os.RemoveAll(tmpUnzip)
	}()

	tmpUnpack, err := unpackRes(ctx, tmpUnzip)
	if err != nil {
		return err
	}
	defer func() {
		_ = os.RemoveAll(tmpUnpack)
	}()

	err = fileutil.MoveAllContent(tmpUnpack, dst)
	if err != nil {
		return err
	}

	return nil
}
func preCheck() error {
	var err error
	err = unpackPreCheck()
	if err != nil {
		return err
	}
	return nil
}
func fetchRes(ctx context.Context, infos []Info) (string, error) {
	tmp, err := os.MkdirTemp("", "arkwaifu-ark-fetch-*")
	if err != nil {
		return "", err
	}
	err = fetch(ctx, infos, tmp)
	if err != nil {
		return "", err
	}
	return tmp, nil
}
func unzipRes(ctx context.Context, tmpFetch string) (string, error) {
	tmp, err := os.MkdirTemp("", "arkwaifu-ark-unzip-*")
	if err != nil {
		return "", err
	}
	err = unzip(ctx, tmpFetch, tmp)
	if err != nil {
		return "", err
	}
	return tmp, nil
}
func unpackRes(ctx context.Context, tmpUnzip string) (string, error) {
	tmp, err := os.MkdirTemp("", "arkwaifu-ark-unpack-*")
	if err != nil {
		return "", err
	}
	err = unpack(ctx, tmpUnzip, tmp)
	if err != nil {
		return "", err
	}
	return tmp, nil
}

// FilterResInfos returns a slice containing all Info matching the specified predicate.
func FilterResInfos(infos []Info, predicate func(i Info) bool) []Info {
	result := make([]Info, 0)
	for _, info := range infos {
		if predicate(info) {
			result = append(result, info)
		}
	}
	return result
}

// FilterResInfosRegexp returns a slice containing all Info matching the specified regexp.
func FilterResInfosRegexp(infos []Info, r []*regexp.Regexp) []Info {
	return FilterResInfos(infos, func(i Info) bool {
		return linq.From(r).
			AnyWith(func(j any) bool {
				return j.(*regexp.Regexp).MatchString(i.Name)
			})
	})
}

// CalculateDifferences returns the infos in "new" but not in "old".
// It checks only "Name" and "MD5" for equality.
// The returned res.Info have the same "resVersion" as the res.Info in "new".
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

func GetFromHGAPI(ctx context.Context, oldVersion string, newVersion string, dest string, filter *regexp.Regexp) error {
	if oldVersion == "" {
		return GetFromHGAPIFully(ctx, newVersion, dest, filter)
	} else {
		return GetFromHGAPIIncrementally(ctx, oldVersion, newVersion, dest, filter)
	}
}
func GetFromHGAPIFully(ctx context.Context, resVersion string, dst string, filters ...*regexp.Regexp) error {
	infos, err := GetResInfos(resVersion)
	if err != nil {
		return errors.WithStack(err)
	}

	infos = FilterResInfosRegexp(infos, filters)

	err = GetRes(ctx, infos, dst)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
func GetFromHGAPIIncrementally(ctx context.Context, oldResVer string, newResVer string, dst string, filters ...*regexp.Regexp) error {
	oldResInfos, err := GetResInfos(oldResVer)
	if err != nil {
		return errors.WithStack(err)
	}
	newResInfos, err := GetResInfos(newResVer)
	if err != nil {
		return errors.WithStack(err)
	}

	oldResInfos = FilterResInfosRegexp(oldResInfos, filters)
	newResInfos = FilterResInfosRegexp(newResInfos, filters)

	infos := CalculateDifferences(newResInfos, oldResInfos)
	err = GetRes(ctx, infos, dst)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
