// Package arkres provides functions to retrieve resources of the game Arknights.
// This package only do fetching, unzipping, unpacking and decrypting resources.
// To do other works, such as parsing BSON or JSON tables or processing graphics assets, use the subpackages.
//
// tools/extractor is required.
package arkres

import (
	"context"
	"github.com/flandiayingman/arkwaifu/internal/pkg/arkres/agdapi"
	"github.com/flandiayingman/arkwaifu/internal/pkg/arkres/hgapi"
	"github.com/pkg/errors"
	"regexp"
)

const (
	DefaultPrefix = "assets/torappu/dynamicassets"
)

func GetLatestResVersion(ctx context.Context) (string, error) {
	agdapiResVersion, err := agdapi.GetLatestResourceVersion(ctx)
	if err != nil {
		return "", err
	}
	return agdapiResVersion.ResourceVersion, nil
}

func GetFromHGAPI(ctx context.Context, resVersion string, dst string, filters ...*regexp.Regexp) error {
	infos, err := hgapi.GetResInfos(resVersion)
	if err != nil {
		return errors.WithStack(err)
	}

	infos = hgapi.FilterResInfosRegexp(infos, filters)

	err = hgapi.GetRes(ctx, infos, dst)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func GetFromHGAPIIncrementally(ctx context.Context, oldResVer string, newResVer string, dst string, filters ...*regexp.Regexp) error {
	oldResInfos, err := hgapi.GetResInfos(oldResVer)
	if err != nil {
		return errors.WithStack(err)
	}
	newResInfos, err := hgapi.GetResInfos(newResVer)
	if err != nil {
		return errors.WithStack(err)
	}

	oldResInfos = hgapi.FilterResInfosRegexp(oldResInfos, filters)
	newResInfos = hgapi.FilterResInfosRegexp(newResInfos, filters)

	infos := hgapi.CalculateDifferences(newResInfos, oldResInfos)
	err = hgapi.GetRes(ctx, infos, dst)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func GetFromAGDAPI(ctx context.Context, resVersion string, dst string, filters ...*regexp.Regexp) error {
	resourceVersion, err := agdapi.GetResourceVersion(ctx, resVersion)
	if err != nil {
		return errors.WithStack(err)
	}
	resourceInfoList, err := agdapi.GetResourceInfoList(ctx, resourceVersion)
	if err != nil {
		return errors.WithStack(err)
	}

	resourceInfoList = agdapi.FilterResInfosRegexp(resourceInfoList, filters)

	err = agdapi.GetRes(ctx, resourceVersion, resourceInfoList, dst)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
