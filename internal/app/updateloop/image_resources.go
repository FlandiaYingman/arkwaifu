package updateloop

import (
	"arkwaifu/internal/app/util/pathutil"
	"context"
	"github.com/chai2010/webp"
	"golang.org/x/image/draw"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
	"image"
	_ "image/png"
	"io/fs"
	"math"
	"os"
	"path/filepath"
)

const concurrency = 16

func createThumbnailOfDir(dirPath string, destDirPath string) error {
	eg, ctx := errgroup.WithContext(context.Background())
	sem := semaphore.NewWeighted(concurrency)
	err := filepath.WalkDir(dirPath, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		imagePath := path
		err = sem.Acquire(ctx, 1)
		if err != nil {
			return err
		}
		eg.Go(func() error {
			defer sem.Release(1)
			rel, _ := filepath.Rel(dirPath, imagePath)
			join := filepath.Join(destDirPath, rel)

			thumbPath := pathutil.ReplaceExt(join, ".webp")
			err := os.MkdirAll(filepath.Dir(thumbPath), 0755)
			if err != nil {
				return err
			}
			err = createThumbnailOf(imagePath, thumbPath)
			if err != nil {
				return err
			}
			return nil
		})
		return nil
	})
	if err != nil {
		_ = eg.Wait()
		return err
	}
	return eg.Wait()
}
func createWebpOfDir(dirPath string, destDirPath string) error {
	eg, ctx := errgroup.WithContext(context.Background())
	sem := semaphore.NewWeighted(concurrency)
	err := filepath.WalkDir(dirPath, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		imagePath := path
		err = sem.Acquire(ctx, 1)
		if err != nil {
			return err
		}
		eg.Go(func() error {
			defer sem.Release(1)
			rel, _ := filepath.Rel(dirPath, imagePath)
			join := filepath.Join(destDirPath, rel)

			thumbPath := pathutil.ReplaceExt(join, ".webp")
			err := os.MkdirAll(filepath.Dir(thumbPath), 0755)
			if err != nil {
				return err
			}
			err = createWebpOf(imagePath, thumbPath)
			if err != nil {
				return err
			}
			return nil
		})
		return nil
	})
	if err != nil {
		_ = eg.Wait()
		return err
	}
	return eg.Wait()
}

func createThumbnailOf(imagePath string, destImagePath string) error {
	imageData, err := decodeImage(imagePath)
	if err != nil {
		return err
	}
	imageData = resizeImage(imageData, 512, 512)
	err = encodeImageToWebp(imageData, destImagePath)
	return err
}
func createWebpOf(imagePath string, destImagePath string) error {
	imageData, err := decodeImage(imagePath)
	if err != nil {
		return err
	}
	err = encodeImageToWebpLossless(imageData, destImagePath)
	return err
}

func decodeImage(imagePath string) (image.Image, error) {
	imageFile, err := os.Open(imagePath)
	if err != nil {
		return nil, err
	}
	defer imageFile.Close()
	imageObj, _, err := image.Decode(imageFile)
	if err != nil {
		return nil, err
	}
	return imageObj, nil
}

func encodeImageToWebp(imageData image.Image, destPath string) error {
	thumbFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer thumbFile.Close()
	err = webp.Encode(thumbFile, imageData, &webp.Options{
		Quality: 85,
	})
	return err
}
func encodeImageToWebpLossless(imageData image.Image, destPath string) error {
	thumbFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer thumbFile.Close()
	err = webp.Encode(thumbFile, imageData, &webp.Options{
		Lossless: true,
	})
	return err
}

func resizeImage(imageData image.Image, maxWidth int, maxHeight int) *image.RGBA {
	width := imageData.Bounds().Dx()
	height := imageData.Bounds().Dy()

	wScale := 1.0
	hScale := 1.0

	if width > maxWidth {
		wScale = float64(maxWidth) / float64(width)
	}
	if height > maxHeight {
		hScale = float64(maxHeight) / float64(height)
	}
	// Since the image must conform to both maxWidth and maxHeight, use the smaller (stricter) scale.
	scale := math.Min(wScale, hScale)
	resizedRect := image.Rect(
		0,
		0,
		int(math.Round(float64(width)*scale)),
		int(math.Round(float64(height)*scale)),
	)
	resizedImage := image.NewRGBA(resizedRect)
	draw.CatmullRom.Scale(resizedImage, resizedRect, imageData, imageData.Bounds(), draw.Over, nil)

	return resizedImage
}
