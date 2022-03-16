package updateloop

import (
	"context"
	"github.com/chai2010/webp"
	"github.com/flandiayingman/arkwaifu/internal/pkg/arkres/asset"
	"github.com/flandiayingman/arkwaifu/internal/pkg/util/fileutil"
	"github.com/flandiayingman/arkwaifu/internal/pkg/util/pathutil"
	"golang.org/x/image/draw"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
	"image"
	_ "image/jpeg" // register jpeg codec
	_ "image/png"  // register png codec
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"regexp"
)

func GetAvgResources(ctx context.Context, oldResVer string, newResVer string, src string, dst string) error {
	assetsRegexp := regexp.MustCompile("^avg/(imgs|bg)")

	tmpDir, err := os.MkdirTemp("", "arkwaifu-updateloop-*")
	if err != nil {
		return err
	}
	defer func() {
		_ = os.RemoveAll(tmpDir)
	}()

	if oldResVer == "" {
		err = asset.Get(ctx, newResVer, tmpDir, assetsRegexp)
		if err != nil {
			return err
		}
	} else {
		err = asset.Update(ctx, oldResVer, newResVer, tmpDir, assetsRegexp)
		if err != nil {
			return err
		}
	}

	rawDir := filepath.Join(dst, "raw")
	err = os.MkdirAll(rawDir, 0755)
	if err != nil {
		return err
	}
	err = processRaw(tmpDir, rawDir)
	if err != nil {
		return err
	}

	thumbnailDir := filepath.Join(dst, "thumbnail")
	err = os.MkdirAll(thumbnailDir, 0755)
	if err != nil {
		return err
	}
	err = processThumbnail(rawDir, thumbnailDir)
	if err != nil {
		return err
	}

	if oldResVer != "" {
		err = fileutil.MoveAllFileContent(src, dst)
		if err != nil {
			return err
		}
	}

	return nil
}

const (
	imagesPath      = "assets/torappu/dynamicassets/avg/images"
	backgroundsPath = "assets/torappu/dynamicassets/avg/backgrounds"
)

func processRaw(src, dst string) error {
	var err error
	imagesPath := filepath.Join(src, imagesPath)
	if _, err := os.Stat(imagesPath); err == nil {
		err = fileutil.MoveAllFileContent(imagesPath, filepath.Join(dst, "images"))
		if err != nil {
			return err
		}
	}
	backgroundsPath := filepath.Join(src, backgroundsPath)
	if _, err := os.Stat(backgroundsPath); err == nil {
		err = fileutil.MoveAllFileContent(backgroundsPath, filepath.Join(dst, "backgrounds"))
		if err != nil {
			return err
		}
	}
	err = fileutil.LowercaseAll(dst)
	return err
}
func processThumbnail(src, dst string) error {
	return createThumbnailOfDir(src, dst)
}

const imgProcConcurrency = 16

func createThumbnailOfDir(dirPath string, destDirPath string) error {
	eg, ctx := errgroup.WithContext(context.Background())
	sem := semaphore.NewWeighted(imgProcConcurrency)
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
	egErr := eg.Wait()
	if egErr != nil {
		return egErr
	}
	if err != nil {
		return err
	}
	return nil
}
func createThumbnailOf(imagePath string, destImagePath string) error {
	imageData, err := decodeImage(imagePath)
	if err != nil {
		return err
	}
	imageData = resizeImage(imageData, 224, 224)
	err = encodeImageToWebp(imageData, destImagePath)
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
