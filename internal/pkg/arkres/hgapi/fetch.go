package hgapi

import (
	"context"
	"path/filepath"

	"github.com/cavaliergopher/grab/v3"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const poolSize = 16
const retryPoolSize = 1

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
	requests = nil
	for response := range responses {
		srcFile := response.Request.URL().String()
		dstFile := filepath.ToSlash(filepath.Clean(response.Filename))

		err := response.Err()
		if err != nil {
			req := response.Request
			req = req.WithContext(ctx)
			requests = append(requests, req)
			log.WithFields(log.Fields{
				"num":   num,
				"total": total,
				"src":   srcFile,
				"dst":   dstFile,
			}).Info("failed to fetch resource; appending it to retry queue")
		} else {
			log.WithFields(log.Fields{
				"num":   num,
				"total": total,
				"src":   srcFile,
				"dst":   dstFile,
			}).Info("fetched resource")
		}
		num++
	}

	responses = grab.DefaultClient.DoBatch(retryPoolSize, requests...)

	num = 0
	total = len(requests)
	for response := range responses {
		srcFile := response.Request.URL().String()
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
		}).Info("fetched resource in retry queue")
		num++
	}
	return nil
}
