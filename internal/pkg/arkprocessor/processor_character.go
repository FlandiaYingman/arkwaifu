package arkprocessor

import (
	"fmt"
	_ "github.com/chai2010/webp"
	"github.com/disintegration/imaging"
	"github.com/flandiayingman/arkwaifu/internal/pkg/arkscanner"
	"github.com/pkg/errors"
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
)

type CharacterArt arkscanner.CharacterArt

type CharacterArtImage struct {
	Image   image.Image
	Art     *CharacterArt
	BodyNum int
	FaceNum int
}

func (i *CharacterArtImage) ID() string {
	return fmt.Sprintf("%s#%d$%d", i.Art.ID, i.FaceNum, i.BodyNum)
}

// ProcessCharacterArt process a character art.
//
// Since picture arts are complicated, the process operation consists of:
//   - Decode all relevant images, including color channel and alpha channel of
//     faces and bodies image.
//   - For each face or body, merge the alpha channels onto the color channels.
//   - For each face, merge it onto its corresponding body.
func (p *Processor) ProcessCharacterArt(art *CharacterArt) ([]CharacterArtImage, error) {
	var result []CharacterArtImage
	for i, body := range art.BodyVariations {
		bodyImage, err := art.decode(p.Root, i+1, 0)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		// Sometimes, there are some character arts that have no actual images given. We'll just skip them.
		// Hope HyperGryph will remember what they've removed and remember to remove them from the stories.
		// Thanks HyperGryph ^^
		if bodyImage == nil {
			continue
		}
		for j, face := range body.FaceVariations {
			faceImage, err := art.decode(p.Root, i+1, j+1)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			// Same reason above.
			// Thanks HyperGryph ^^
			if faceImage == nil {
				continue
			}
			if !face.WholeBody {
				faceImage, err = mergeCharacterFace(bodyImage, faceImage, body.FaceRectangle)
				if err != nil {
					return nil, errors.WithStack(err)
				}
			}
			result = append(result, CharacterArtImage{
				Image:   faceImage,
				Art:     art,
				BodyNum: i + 1,
				FaceNum: j + 1,
			})
		}
	}
	return result, nil
}

func (a *CharacterArt) decode(root string, bodyNum, faceNum int) (image.Image, error) {
	var filePath, filePathAlpha string
	if faceNum > 0 {
		filePath = filepath.Join(root, (*arkscanner.CharacterArt)(a).FacePath(bodyNum, faceNum))
		filePathAlpha = filepath.Join(root, (*arkscanner.CharacterArt)(a).FacePathAlpha(bodyNum, faceNum))
	} else {
		filePath = filepath.Join(root, (*arkscanner.CharacterArt)(a).BodyPath(bodyNum))
		filePathAlpha = filepath.Join(root, (*arkscanner.CharacterArt)(a).BodyPathAlpha(bodyNum))
	}
	if filePath == root {
		return nil, nil
	}

	var img, imgAlpha image.Image
	var err error
	img, err = decode(filePath)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if filePathAlpha != root {
		imgAlpha, err = decode(filePathAlpha)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}

	if imgAlpha != nil {
		img, err = mergeAlphaChannel(img, imgAlpha)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}

	return img, nil
}

func decode(filePath string) (image.Image, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer func() { _ = file.Close() }()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return img, nil
}
func mergeAlphaChannel(base image.Image, alpha image.Image) (*image.NRGBA, error) {
	// TODO: DrawMask instead
	if base.Bounds() != alpha.Bounds() {
		alpha = imaging.Resize(alpha, base.Bounds().Dx(), base.Bounds().Dy(), imaging.Lanczos)
	}
	canvas := image.NewNRGBA(base.Bounds())
	for i := canvas.Bounds().Min.X; i < canvas.Bounds().Max.X; i++ {
		for j := canvas.Bounds().Min.Y; j < canvas.Bounds().Max.Y; j++ {
			baseColor := color.NRGBAModel.Convert(base.At(i, j)).(color.NRGBA)
			alphaColor := color.GrayModel.Convert(alpha.At(i, j)).(color.Gray)
			baseColor.A = alphaColor.Y
			canvas.SetNRGBA(i, j, baseColor)
		}
	}
	return canvas, nil
}

func mergeCharacterFace(body image.Image, face image.Image, faceRect image.Rectangle) (*image.NRGBA, error) {
	if body == nil || face == nil {
		return nil, errors.Errorf("the images are nil! ")
	}
	if !faceRect.In(body.Bounds()) {
		return nil, errors.Errorf("merge character face: face rectangle %v is not in the body's bounds %v", faceRect, body.Bounds())
	}
	face = imaging.Resize(face, faceRect.Dx(), faceRect.Dy(), imaging.Lanczos)

	canvas := image.NewNRGBA(body.Bounds())
	draw.Draw(canvas, body.Bounds(), body, image.Point{}, draw.Over)
	draw.Draw(canvas, faceRect, face, image.Point{}, draw.Over)
	return canvas, nil
}
