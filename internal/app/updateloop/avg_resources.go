package updateloop

import (
	"arkwaifu/internal/app/util/fileutil"
	"arkwaifu/internal/app/util/pathutil"
	"arkwaifu/internal/pkg/arkres/resource"
	"context"
	"github.com/chai2010/webp"
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

func GetAvgResources(resVersion string, dest string) error {
	infos, err := resource.GetResInfos(resVersion)
	if err != nil {
		return err
	}

	tmpDir, err := os.MkdirTemp("", "arkwaifu-updateloop-avg_resources-*")
	if err != nil {
		return err
	}
	infos = resource.FilterResInfosRegexp(infos, regexp.MustCompile("^avg/"))
	err = resource.GetRes(infos, tmpDir)
	if err != nil {
		return err
	}

	rawDir := filepath.Join(dest, "raw")
	err = os.MkdirAll(rawDir, 0755)
	if err != nil {
		return err
	}
	err = fileutil.MoveAllFileContent(filepath.Join(tmpDir, "assets/torappu/dynamicassets/avg/images"), filepath.Join(rawDir, "images"))
	if err != nil {
		return err
	}
	err = fileutil.MoveAllFileContent(filepath.Join(tmpDir, "assets/torappu/dynamicassets/avg/backgrounds"), filepath.Join(rawDir, "backgrounds"))
	if err != nil {
		return err
	}
	err = os.RemoveAll(tmpDir)
	if err != nil {
		return err
	}

	thumbnailDir := filepath.Join(dest, "thumbnail")
	err = os.MkdirAll(thumbnailDir, 0755)
	if err != nil {
		return err
	}
	err = createThumbnailOfDir(rawDir, thumbnailDir)
	if err != nil {
		return err
	}

	return nil
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
	imageData = resizeImage(imageData, 512, 512)
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
