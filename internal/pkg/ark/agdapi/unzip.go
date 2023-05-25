package agdapi

import (
	"context"
	"github.com/flandiayingman/arkwaifu/internal/pkg/ark"
	"github.com/flandiayingman/arkwaifu/internal/pkg/cols"
	"github.com/flandiayingman/arkwaifu/internal/pkg/util/cloze"
	"github.com/flandiayingman/arkwaifu/internal/pkg/util/fileutil"
	"github.com/flandiayingman/arkwaifu/internal/pkg/util/pathutil"
	"github.com/mholt/archiver/v4"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func unzipZipball(ctx context.Context, zipballPath string, resourceInfoList []ResourceInfo, resourceVersion ResourceVersion, destination string) error {
	// Since the zipball contains a directory directly named by the stem of the zipball,
	// the prefix is gotten for generating the file path inside the zipball.
	filePrefix := pathutil.Stem(zipballPath)
	fileList := cols.Map(resourceInfoList, func(info ResourceInfo) string {
		return path.Join(filePrefix, ark.LanguageCodeUnderscore(resourceVersion.GameServer), info.Name)
	})

	reader, err := os.Open(zipballPath)
	if err != nil {
		return err
	}
	defer cloze.IgnoreErr(reader)

	zip := archiver.Zip{}
	return zip.Extract(ctx, reader, fileList, func(ctx context.Context, file archiver.File) error {
		resourceName := path.Join(strings.Split(file.NameInArchive, "/")[2:]...)
		destinationPath := filepath.Join(destination, DefaultPrefix, resourceName)

		reader, err := file.Open()
		if err != nil {
			return err
		}
		defer cloze.IgnoreErr(reader)

		err = fileutil.MkFileFromReader(destinationPath, reader)
		if err != nil {
			return err
		}
		return nil
	})
}
