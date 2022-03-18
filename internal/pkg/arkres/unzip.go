package arkres

import (
	"context"
	"github.com/mholt/archiver/v4"
	log "github.com/sirupsen/logrus"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
)

func unzip(ctx context.Context, src string, dst string) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if err := ctx.Err(); err != nil {
			return err
		}

		err = unzipFile(ctx, path, dst)
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

		log.WithFields(log.Fields{
			"src": src,
			"dst": dst,
		}).Info("unzipped resource")
		return nil
	})
	return err
}
