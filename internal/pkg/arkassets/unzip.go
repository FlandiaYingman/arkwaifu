package arkassets

import (
	"context"
	"github.com/mholt/archiver/v4"
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
		return "", err
	}

	return tempDir, filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if err := ctx.Err(); err != nil {
			return err
		}

		err = unzipFile(ctx, path, tempDir)
		if err != nil {
			return err
		}
		return nil
	})
}

func unzipFile(ctx context.Context, src string, dst string) error {
	zip := archiver.Zip{}

	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	err = zip.Extract(ctx, srcFile, nil, func(ctx context.Context, f archiver.File) error {
		if err := ctx.Err(); err != nil {
			return err
		}

		compFile, err := f.Open()
		if err != nil {
			return err
		}
		defer compFile.Close()

		compFileName := f.NameInArchive

		src := filepath.ToSlash(filepath.Clean(src))
		dst := filepath.ToSlash(filepath.Clean(path.Join(dst, compFileName)))
		err = os.MkdirAll(filepath.Dir(dst), os.ModePerm)
		if err != nil {
			return err
		}

		dstFile, err := os.Create(dst)
		if err != nil {
			return err
		}
		defer dstFile.Close()

		_, err = io.Copy(dstFile, compFile)
		if err != nil {
			return err
		}

		log.Info().
			Str("src", src).
			Str("dst", dst).
			Msg("Unzipped resource src to dst.")

		return nil
	})
	return err
}
