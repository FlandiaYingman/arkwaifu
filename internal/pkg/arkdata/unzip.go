package arkdata

import (
	"context"
	"github.com/flandiayingman/arkwaifu/internal/pkg/ark"
	"github.com/flandiayingman/arkwaifu/internal/pkg/cols"
	"github.com/flandiayingman/arkwaifu/internal/pkg/util/cloze"
	"github.com/flandiayingman/arkwaifu/internal/pkg/util/fileutil"
	"github.com/flandiayingman/arkwaifu/internal/pkg/util/pathutil"
	"github.com/mholt/archiver/v4"
	"github.com/pkg/errors"
	"os"
	"path"
	"path/filepath"
)

func unzip(ctx context.Context, zipball string, patterns []string, server ark.Server) (string, error) {
	temp, err := os.MkdirTemp("", "arkdata_unzip")
	if err != nil {
		return "", errors.WithStack(err)
	}

	// Since the zipball contains a directory directly named by the stem of the zipball,
	// the prefix is gotten for generating the file path inside the zipball.
	patterns = cols.Map(patterns, func(pattern string) string {
		return path.Join(pathutil.Stem(zipball), ark.LanguageCodeUnderscore(server), pattern)
	})

	reader, err := os.Open(zipball)
	if err != nil {
		return "", errors.WithStack(err)
	}
	defer cloze.IgnoreErr(reader)

	zip := archiver.Zip{}
	return temp, zip.Extract(ctx, reader, nil, func(ctx context.Context, file archiver.File) error {
		if file.IsDir() {
			return nil
		}

		match, err := pathutil.MatchAny(patterns, file.NameInArchive)
		if err != nil {
			return errors.WithStack(err)
		}
		if !match {
			return nil
		}

		fileName := path.Join(pathutil.Splits(file.NameInArchive)[2:]...)
		filePath := filepath.Join(temp, DefaultPrefix, fileName)

		reader, err := file.Open()
		if err != nil {
			return errors.WithStack(err)
		}

		err = fileutil.MkFileFromReader(filePath, reader)
		if err != nil {
			return errors.WithStack(err)
		}

		return nil
	})
}
