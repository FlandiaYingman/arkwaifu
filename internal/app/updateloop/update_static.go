package updateloop

import (
	"context"
	"image"
	"math"
	"os"
	"path/filepath"
	"runtime"

	"github.com/chai2010/webp"
	"github.com/flandiayingman/arkwaifu/internal/pkg/arkres"
	"github.com/flandiayingman/arkwaifu/internal/pkg/arkres/arkavg"
	"github.com/flandiayingman/arkwaifu/internal/pkg/util/fileutil"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/image/draw"
	"golang.org/x/sync/errgroup"
)

// updateStatic processes the resources to static files.
// The resources are from the corresponding resource directory,
// and the static files are stored in the corresponding static directory.
//
// It is skipped if the corresponding static directory already exists.
//
// Currently, the static files consist of the following:
//
// The image assets (img), processed into the WebP format from the raw resources;
// The thumbnail assets (timg), processed into the WebP format from the raw resources;
// ...
func (s *Service) updateStatic(ctx context.Context, resVer ResVersion) error {
	resDir := s.ResDir(resVer)
	staticDir := s.StaticDir(resVer)

	exists, err := fileutil.Exists(staticDir)
	if err != nil {
		return errors.Wrapf(err, "cannot check if staticDir exists: %s", staticDir)
	}
	if exists && !s.ForceUpdate {
		log.WithFields(log.Fields{
			"staticDir": staticDir,
		}).Info("update static: staticDir already exists; skipping")
		return nil
	}

	err = updateStatic(ctx, resDir, staticDir)
	if err != nil {
		return errors.Wrapf(err, "cannot update statics: %s", staticDir)
	}

	log.WithFields(log.Fields{
		"resDir":    resDir,
		"staticDir": staticDir,
	}).Info("processed statics from resDir to staticDir")
	return nil
}

// updateStatic processes the resources into static files.
//
// See: Service.updateStatic.
func updateStatic(ctx context.Context, newResDir string, newStaticDir string) error {
	assets, err := arkavg.GetAssets(newResDir, arkres.DefaultPrefix)
	if err != nil {
		// Skip if no resources to process
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	eg, ctx := errgroup.WithContext(ctx)
	lock := make(chan struct{}, runtime.NumCPU())
	for _, asset := range assets {
		lock <- struct{}{}
		asset := asset
		eg.Go(func() error {
			err := processStatic(asset, newResDir, newStaticDir)
			<-lock
			return err
		})
	}
	err = eg.Wait()
	if err != nil {
		return err
	}

	return nil
}

func processStatic(asset arkavg.Asset, resDir string, staticDir string) error {
	assetPath := asset.FilePath(resDir, arkres.DefaultPrefix)
	assetFile, err := os.Open(assetPath)
	if err != nil {
		return errors.Wrapf(err, "cannot open asset: %s", assetPath)
	}

	img, _, err := image.Decode(assetFile)
	if err != nil {
		return errors.Wrapf(err, "cannot decode asset: %s", assetPath)
	}

	err = encodeImg(asset, img, staticDir)
	if err != nil {
		return err
	}
	err = encodeTimg(asset, img, staticDir)
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"asset":     asset,
		"staticDir": staticDir,
	}).Info("processed static file from asset into static dir")

	return nil
}
func encodeImg(asset arkavg.Asset, img image.Image, staticDir string) error {
	imgStaticDir := filepath.Join(staticDir, "img")

	dstPath := filepath.Join(imgStaticDir, string(asset.Kind), asset.Name+".webp")
	dstFile, err := fileutil.MkFile(dstPath)
	if err != nil {
		return errors.Wrapf(err, "cannot create dstFile: %s", dstPath)
	}
	defer func() { _ = dstFile.Close() }()

	return webp.Encode(dstFile, img, &webp.Options{
		Lossless: true,
		Exact:    true,
	})
}
func encodeTimg(asset arkavg.Asset, img image.Image, staticDir string) error {
	timgStaticDir := filepath.Join(staticDir, "timg")

	dstPath := filepath.Join(timgStaticDir, string(asset.Kind), asset.Name+".webp")
	dstFile, err := fileutil.MkFile(dstPath)
	if err != nil {
		return errors.Wrapf(err, "cannot create dstFile: %s", dstPath)
	}
	defer func() { _ = dstFile.Close() }()

	img = resize(img, 320, 1280)
	return webp.Encode(dstFile, img, &webp.Options{
		Quality: 85,
	})
}
func resize(img image.Image, maxW int, maxH int) image.Image {
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	wScale := 1.0
	hScale := 1.0
	if width > maxW {
		wScale = float64(maxW) / float64(width)
	}
	if height > maxH {
		hScale = float64(maxH) / float64(height)
	}
	// Since the image must conform to both maxW and maxH, use the smaller (stricter) scale.
	scale := math.Min(wScale, hScale)
	resultRect := image.Rect(0, 0,
		int(math.Round(float64(width)*scale)),
		int(math.Round(float64(height)*scale)),
	)
	result := image.NewNRGBA(resultRect)
	draw.ApproxBiLinear.Scale(result, resultRect, img, img.Bounds(), draw.Over, nil)

	return result
}
