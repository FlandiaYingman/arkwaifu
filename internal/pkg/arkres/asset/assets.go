package asset

import (
	"context"
	"github.com/ahmetb/go-linq/v3"
	"github.com/flandiayingman/arkwaifu/internal/pkg/arkres/res"
	"github.com/pkg/errors"
	"regexp"
)

func Get(ctx context.Context, resVersion string, dst string, filter *regexp.Regexp) error {
	infos, err := res.GetResInfos(resVersion)
	if err != nil {
		return errors.WithStack(err)
	}

	infos = res.FilterResInfosRegexp(infos, filter)

	err = res.GetRes(ctx, infos, dst)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func Update(ctx context.Context, oldResVer string, newResVer string, dst string, filter *regexp.Regexp) error {
	oldResInfos, err := res.GetResInfos(oldResVer)
	if err != nil {
		return errors.WithStack(err)
	}
	newResInfos, err := res.GetResInfos(newResVer)
	if err != nil {
		return errors.WithStack(err)
	}

	oldResInfos = res.FilterResInfosRegexp(oldResInfos, filter)
	newResInfos = res.FilterResInfosRegexp(newResInfos, filter)

	infos := calcDiff(newResInfos, oldResInfos)
	err = res.GetRes(ctx, infos, dst)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func GetLatestResVersion() (string, error) {
	return res.GetResVersion()
}

// calcDiff returns the infos in "new" but not in "old".
// It checks only "Name" and "MD5" for equality.
// The returned res.Info have the same "resVersion" as the res.Info in "new".
func calcDiff(new []res.Info, old []res.Info) []res.Info {
	convert := func(i interface{}) interface{} {
		info := i.(res.Info)
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
	var result []res.Info
	newQuery.
		ExceptBy(oldQuery, convert).
		ToSlice(&result)
	return result
}
