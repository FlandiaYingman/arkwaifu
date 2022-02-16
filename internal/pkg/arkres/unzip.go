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

func unzip(src string, dst string) error {
	zip := archiver.Zip{}

	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer closeFile(srcFile)

	err = zip.Extract(context.Background(), srcFile, nil, func(ctx context.Context, f archiver.File) error {
		compFile, err := f.Open()
		if err != nil {
			return err
		}
		defer closeFile(compFile)

		compressedFileName := f.NameInArchive

		src := filepath.ToSlash(filepath.Clean(src))
		dst := filepath.ToSlash(filepath.Clean(path.Join(dst, compressedFileName)))
		err = os.MkdirAll(filepath.Dir(dst), os.ModePerm)
		if err != nil {
			return err
		}

		dstFile, err := os.Create(dst)
		if err != nil {
			return err
		}
		defer closeFile(dstFile)

		_, err = io.Copy(dstFile, compFile)
		if err != nil {
			return err
		}

		log.WithFields(log.Fields{
			"src": src,
			"dst": dst,
		}).Info("Resource unzipped.")
		return nil
	})
	return err
}

func closeFile(file io.ReadCloser) {
	err := file.Close()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Warn("err while closing")
	}
}

func unzipResources(src string, dst string) error {
	err := filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		err = unzip(path, dst)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
