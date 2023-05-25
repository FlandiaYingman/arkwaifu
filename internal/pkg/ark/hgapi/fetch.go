package hgapi

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/cavaliergopher/grab/v3"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
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

		log := log.With().
			Str("src", srcFile).
			Str("dst", dstFile).
			Int("num", num).
			Int("total", total).
			Str("percent", fmt.Sprintf("%.2f%", float64(num)/float64(total))).
			Logger()

		err := response.Err()
		if err != nil {
			req := response.Request
			req = req.WithContext(ctx)
			requests = append(requests, req)
			log.Info().Msg("Failed to fetch resource; appending to retry queue")
		} else {
			log.Info().Msg("Fetched resource.")
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
		log.Info().
			Str("src", srcFile).
			Str("dst", dstFile).
			Int("num", num).
			Int("total", total).
			Str("percent", fmt.Sprintf("%.2f%", float64(num)/float64(total))).
			Msg("Fetched resource in retry queue.")
		num++
	}
	return nil
}
