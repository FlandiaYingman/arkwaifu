package arkres

import (
	"context"
	"encoding/base64"
	"os"
	"regexp"
	"strings"

	"github.com/ahmetb/go-linq/v3"
	"github.com/caarlos0/env/v6"
	"github.com/flandiayingman/arkwaifu/internal/pkg/util/fileutil"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func init() {
	conf := struct {
		ExtractorLocation *string `env:"ARKRES_EXTRACTOR_LOCATION"`
		ChatMask          *string `env:"ARKRES_CHAT_MASK"`
	}{
		ExtractorLocation: nil,
		ChatMask:          nil,
	}
	err := env.Parse(&conf)
	if err != nil {
		log.WithError(err).Warn("parsing environment variables")
		return
	}

	if conf.ExtractorLocation != nil {
		extractorLocation = *conf.ExtractorLocation
	}
	if conf.ChatMask != nil {
		bytesChatMask, err := base64.StdEncoding.DecodeString(*conf.ChatMask)
		if err != nil {
			log.
				WithError(err).
				WithField("conf", conf).
				Warn("setting chat mask")
			return
		}
		SetChatMask(bytesChatMask)
	}
}

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

	tmpDecrypt, err := decryptRes(ctx, tmpUnpack)
	if err != nil {
		return err
	}
	defer func() {
		_ = os.RemoveAll(tmpDecrypt)
	}()

	err = fileutil.MoveAllContent(tmpUnpack, dst)
	if err != nil {
		return err
	}
	err = fileutil.MoveAllContent(tmpDecrypt, dst)
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
	err = decryptPreCheck()
	if err != nil {
		return err
	}
	return nil
}
func fetchRes(ctx context.Context, infos []Info) (string, error) {
	tmp, err := os.MkdirTemp("", "arkwaifu-arkres-fetch-*")
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
	tmp, err := os.MkdirTemp("", "arkwaifu-arkres-unzip-*")
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
	tmp, err := os.MkdirTemp("", "arkwaifu-arkres-unpack-*")
	if err != nil {
		return "", err
	}
	err = unpack(ctx, tmpUnzip, tmp)
	if err != nil {
		return "", err
	}
	return tmp, nil
}
func decryptRes(ctx context.Context, tmpUnpack string) (string, error) {
	tmp, err := os.MkdirTemp("", "arkwaifu-arkres-decrypt-*")
	if err != nil {
		return "", err
	}
	err = decrypt(ctx, tmpUnpack, tmp)
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
