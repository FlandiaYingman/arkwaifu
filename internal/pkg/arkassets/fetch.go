package arkassets

import (
	"context"
	"fmt"
	"github.com/flandiayingman/arkwaifu/internal/pkg/cols"
	"os"
	"path/filepath"

	"github.com/cavaliergopher/grab/v3"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

const fetchWorkers = 16
const retryWorkers = 1

func fetch(ctx context.Context, infoList []Info) (string, error) {
	tempDir, err := os.MkdirTemp("", "arkassets_fetch-*")
	if err != nil {
		return "", errors.WithStack(err)
	}

	ctx, cancelCtx := context.WithCancel(ctx)
	defer cancelCtx()

	requests, err := cols.MapErr(infoList, func(i Info) (*grab.Request, error) {
		request, err := grab.NewRequest(tempDir, i.Url())
		if err != nil {
			return nil, errors.WithStack(err)
		}
		request = request.WithContext(ctx)
		return request, nil
	})
	responses := grab.DefaultClient.DoBatch(fetchWorkers, requests...)

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
			Str("progress", fmt.Sprintf("%.3f", float64(num)/float64(total))).
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

	responses = grab.DefaultClient.DoBatch(retryWorkers, requests...)

	num = 0
	total = len(requests)
	for response := range responses {
		srcFile := response.Request.URL().String()
		dstFile := filepath.ToSlash(filepath.Clean(response.Filename))

		err := response.Err()
		if err != nil {
			return "", errors.Wrapf(err, "srcFile：%v; dstFile: %v", srcFile, dstFile)
		}
		log.Info().
			Str("src", srcFile).
			Str("dst", dstFile).
			Int("num", num).
			Int("total", total).
			Str("progress", fmt.Sprintf("%.3f", float64(num)/float64(total))).
			Msg("Fetched resource in retry queue.")
		num++
	}
	return tempDir, nil
}
