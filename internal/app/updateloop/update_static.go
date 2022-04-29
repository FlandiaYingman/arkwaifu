package updateloop

import (
	"context"
	"image"
	"os"
	"path/filepath"
	"runtime"

	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
	"github.com/flandiayingman/arkwaifu/internal/pkg/arkres"
	"github.com/flandiayingman/arkwaifu/internal/pkg/arkres/arkavg"
	"github.com/flandiayingman/arkwaifu/internal/pkg/util/fileutil"
	"github.com/flandiayingman/arkwaifu/internal/pkg/util/pathutil"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

// updatePicStatic processes the resources to static files.
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

	err = updatePicStatic(ctx, resDir, staticDir)
	if err != nil {
		return errors.Wrapf(err, "cannot update statics: %s", staticDir)
	}
	err = updateCharStatic(ctx, resDir, staticDir)
	if err != nil {
		return errors.Wrapf(err, "cannot update statics: %s", staticDir)
	}

	log.WithFields(log.Fields{
		"resDir":    resDir,
		"staticDir": staticDir,
	}).Info("processed statics from resDir to staticDir")
	return nil
}

func updatePicStatic(ctx context.Context, newResDir string, newStaticDir string) error {
	assets, err := arkavg.ScanForPicAssets(newResDir, arkres.DefaultPrefix)
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
			err := procPicAsset(asset, newResDir, newStaticDir)
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
func procPicAsset(asset arkavg.Asset, resDir string, staticDir string) error {
	assetPath := asset.FilePath(resDir, arkres.DefaultPrefix)
	assetFile, err := os.Open(assetPath)
	if err != nil {
		return errors.Wrapf(err, "cannot open asset: %s", assetPath)
	}

	img, _, err := image.Decode(assetFile)
	if err != nil {
		return errors.Wrapf(err, "cannot decode asset: %s", assetPath)
	}

	err = saveIMG(asset, img, staticDir)
	if err != nil {
		return err
	}
	err = saveTIMG(asset, img, staticDir)
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"asset":     asset,
		"staticDir": staticDir,
	}).Info("processed static file from asset into static dir")

	return nil
}

func updateCharStatic(ctx context.Context, resDir, staticDir string) error {
	assets, err := arkavg.ScanForCharAssets(resDir, arkres.DefaultPrefix)
	if errors.Is(err, os.ErrNotExist) {
		// Skip if no assets to process
		return nil
	}
	if err != nil {
		return err
	}
	eg, ctx := errgroup.WithContext(ctx)
	lock := make(chan struct{}, runtime.NumCPU())
	for _, asset := range assets {
		lock <- struct{}{}
		asset := asset
		eg.Go(func() error {
			err := procCharAsset(asset, resDir, staticDir)
			<-lock
			return err
		})
	}
	return eg.Wait()
}
func procCharAsset(ca arkavg.CharAsset, resDir, staticDir string) error {
	for hn, h := range ca.Hubs {
		merges, mergesAlpha, err := h.MergeFace(resDir, arkres.DefaultPrefix)
		if err != nil {
			return err
		}
		for fn := range merges {
			asset := arkavg.Asset{
				Name: ca.AssetName(hn, fn),
				Kind: arkavg.KindCharacter,
			}
			err = saveIMG(asset, merges[fn], staticDir)
			if err != nil {
				return err
			}
			err = saveAlpha(asset, mergesAlpha[fn], staticDir)
			if err != nil {
				return err
			}
			err = saveTIMG(asset, arkavg.MergeAlpha(merges[fn], mergesAlpha[fn]), staticDir)
			if err != nil {
				return err
			}
		}
	}
	log.WithFields(log.Fields{
		"charAsset": ca,
		"staticDir": staticDir,
	}).Info("processed static file from char asset into static dir")
	return nil
}

func saveIMG(asset arkavg.Asset, img image.Image, staticDir string) error {
	imgStaticDir := filepath.Join(staticDir, "img")

	dstPath := filepath.Join(imgStaticDir, string(asset.Kind), pathutil.ReplaceAllExt(asset.Name, ".webp"))
	dstFile, err := fileutil.MkFile(dstPath)
	if err != nil {
		return errors.Wrapf(err, "cannot create dstFile: %s", dstPath)
	}
	defer func() { _ = dstFile.Close() }()
	return webp.Encode(dstFile, img, &webp.Options{
		Lossless: true,
	})
}
func saveTIMG(asset arkavg.Asset, img image.Image, staticDir string) error {
	timgStaticDir := filepath.Join(staticDir, "timg")

	dstPath := filepath.Join(timgStaticDir, string(asset.Kind), pathutil.ReplaceAllExt(asset.Name, ".webp"))
	dstFile, err := fileutil.MkFile(dstPath)
	if err != nil {
		return errors.Wrapf(err, "cannot create dstFile: %s", dstPath)
	}
	defer func() { _ = dstFile.Close() }()
	img = imaging.Fit(img, 320, 1280, imaging.Linear)
	return webp.Encode(dstFile, img, &webp.Options{
		Quality: 85,
	})
}
func saveAlpha(asset arkavg.Asset, img image.Image, staticDir string) error {
	imgStaticDir := filepath.Join(staticDir, "alpha")

	dstPath := filepath.Join(imgStaticDir, string(asset.Kind), pathutil.ReplaceAllExt(asset.Name, ".webp"))
	dstFile, err := fileutil.MkFile(dstPath)
	if err != nil {
		return errors.Wrapf(err, "cannot create dstFile: %s", dstPath)
	}
	defer func() { _ = dstFile.Close() }()
	return webp.Encode(dstFile, img, &webp.Options{
		Lossless: true,
	})
}
