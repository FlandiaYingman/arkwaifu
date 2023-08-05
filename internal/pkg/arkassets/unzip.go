package arkassets

import (
	"context"
	"github.com/mholt/archiver/v4"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
)

func unzip(ctx context.Context, src string) (string, error) {
	tempDir, err := os.MkdirTemp("", "arkassets_unzip-*")
	if err != nil {
		return "", errors.WithStack(err)
	}

	return tempDir, filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return errors.WithStack(err)
		}
		if d.IsDir() {
			return nil
		}
		if err := ctx.Err(); err != nil {
			return errors.WithStack(err)
		}

		err = unzipFile(ctx, path, tempDir)
		if err != nil {
			return errors.WithStack(err)
		}
		return nil
	})
}

func unzipFile(ctx context.Context, src string, dst string) error {
	zip := archiver.Zip{}

	srcFile, err := os.Open(src)
	if err != nil {
		return errors.WithStack(err)
	}
	defer srcFile.Close()

	err = zip.Extract(ctx, srcFile, nil, func(ctx context.Context, f archiver.File) error {
		if err := ctx.Err(); err != nil {
			return errors.WithStack(err)
		}

		compFile, err := f.Open()
		if err != nil {
			return errors.WithStack(err)
		}
		defer compFile.Close()

		compFileName := f.NameInArchive

		src := filepath.ToSlash(filepath.Clean(src))
		dst := filepath.ToSlash(filepath.Clean(path.Join(dst, compFileName)))
		err = os.MkdirAll(filepath.Dir(dst), os.ModePerm)
		if err != nil {
			return errors.WithStack(err)
		}

		dstFile, err := os.Create(dst)
		if err != nil {
			return errors.WithStack(err)
		}
		defer dstFile.Close()

		_, err = io.Copy(dstFile, compFile)
		if err != nil {
			return errors.WithStack(err)
		}

		log.Info().
			Str("src", src).
			Str("dst", dst).
			Msg("Unzipped resource src to dst.")

		return nil
	})
	return errors.WithStack(err)
}
