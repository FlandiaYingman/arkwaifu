package updateloop

import (
	"context"
	"github.com/chai2010/webp"
	"github.com/flandiayingman/arkwaifu/internal/pkg/arkres"
	"github.com/flandiayingman/arkwaifu/internal/pkg/arkres/arkavg"
	"github.com/pkg/errors"
	"golang.org/x/image/draw"
	"golang.org/x/sync/errgroup"
	"image"
	"math"
	"os"
	"path/filepath"
	"runtime"
)

// updateStatics reads resources from resLoc and process them into the format of static assets.
// As follows:
//
// Converting the assets from the format of png into webp.
// Converting the assets from original into thumbnails, also in webp format.
// ...
func updateStatics(ctx context.Context, newResDir string, newStaticDir string) error {
	assets, err := arkavg.GetAssets(newResDir, arkres.DefaultPrefix)
	if err != nil {
		return err
	}

	eg, ctx := errgroup.WithContext(ctx)
	concurrency := make(chan struct{}, runtime.NumCPU()+1)
	for _, asset := range assets {
		concurrency <- struct{}{}
		asset, newResDir, newStaticDir := asset, newResDir, newStaticDir
		eg.Go(func() error {
			err := procAsset(asset, newResDir, newStaticDir)
			<-concurrency
			return err
		})
	}
	err = eg.Wait()
	if err != nil {
		return err
	}

	return nil
}

func procAsset(asset arkavg.Asset, resDir string, staticDir string) error {
	assetFile, err := asset.Open(resDir, arkres.DefaultPrefix)
	if err != nil {
		return errors.WithStack(err)
	}
	defer func() { _ = assetFile.Close() }()

	assetImg, _, err := image.Decode(assetFile)
	if err != nil {
		return errors.WithStack(err)
	}

	imgStaticDir := filepath.Join(staticDir, "img")
	timgStaticDir := filepath.Join(staticDir, "timg")

	err = encodeImg(asset, assetImg, imgStaticDir)
	if err != nil {
		return err
	}
	err = encodeTimg(asset, assetImg, timgStaticDir)
	if err != nil {
		return err
	}

	return nil
}

func encodeImg(asset arkavg.Asset, img image.Image, imgDir string) error {
	dst := filepath.Join(imgDir, string(asset.Kind), asset.ID+".webp")
	err := os.MkdirAll(filepath.Dir(dst), 0755)
	if err != nil {
		return errors.WithStack(err)
	}

	dstFile, err := os.Create(dst)
	if err != nil {
		return errors.WithStack(err)
	}
	defer func() { _ = dstFile.Close() }()
	return webp.Encode(dstFile, img, &webp.Options{
		Lossless: true,
		Exact:    true,
	})
}

func encodeTimg(asset arkavg.Asset, img image.Image, imgDir string) error {
	dst := filepath.Join(imgDir, string(asset.Kind), asset.ID+".webp")
	err := os.MkdirAll(filepath.Dir(dst), 0755)
	if err != nil {
		return errors.WithStack(err)
	}

	img = resize(img, 320, 1024)

	dstFile, err := os.Create(dst)
	if err != nil {
		return errors.WithStack(err)
	}
	defer func() { _ = dstFile.Close() }()
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
