package res

import (
	"context"
	"github.com/cavaliergopher/grab/v3"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"path"
	"path/filepath"
)

const poolSize = 24

// fetch fetches the resources specified by srcs into dst.
// dst must exist.
func fetch(ctx context.Context, srcs []Info, dst string) error {
	ctx, cancelCtx := context.WithCancel(ctx)
	defer cancelCtx()

	requests := make([]*grab.Request, len(srcs))
	for i, src := range srcs {
		request, err := grab.NewRequest(dst, src.CreateURL())
		if err != nil {
			return err
		}

		requests[i] = request.WithContext(ctx)
	}

	responses := grab.DefaultClient.DoBatch(poolSize, requests...)

	num := 0
	total := len(requests)
	for response := range responses {
		srcFile := path.Clean(response.Request.URL().String())
		dstFile := filepath.ToSlash(filepath.Clean(response.Filename))

		err := response.Err()
		if err != nil {
			return errors.Wrapf(err, "srcFile：%v; dstFile: %v", srcFile, dstFile)
		}

		log.WithFields(log.Fields{
			"num":   num,
			"total": total,
			"src":   srcFile,
			"dst":   dstFile,
		}).Info("fetched resource")
		num++
	}
	return nil
}
