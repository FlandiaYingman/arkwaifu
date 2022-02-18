package resource

import (
	"context"
	"github.com/cavaliergopher/grab/v3"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"path/filepath"
)

const poolSize = 16

func downloadResources(resourceUrls []string, dst string) error {
	err := os.MkdirAll(dst, os.ModePerm)
	if err != nil {
		return err
	}

	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	requests := make([]*grab.Request, len(resourceUrls))
	for i, url := range resourceUrls {
		request, err := grab.NewRequest(dst, url)
		if err != nil {
			return err
		}

		request = request.WithContext(ctx)
		requests[i] = request
	}

	client := grab.NewClient()
	responses := client.DoBatch(poolSize, requests...)

	current := 0
	total := len(requests)
	for response := range responses {
		err := response.Err()
		if err != nil {
			return err
		}
		current++

		url := path.Clean(response.Request.URL().String())
		filename := filepath.ToSlash(filepath.Clean(response.Filename))
		log.WithFields(log.Fields{
			"current":  current,
			"total":    total,
			"url":      url,
			"filename": filename,
		}).Info("Resource downloaded.")
	}
	return nil
}
