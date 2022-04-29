package arkavg

import (
	"image"
	"image/color"
	"image/draw"
	_ "image/png"
	"os"

	"github.com/disintegration/imaging"
)

type pos struct {
	X, Y float64
}

func (h *CharAssetSpriteHub) MergeFace(resDir, prefix string) ([]image.Image, []image.Image, error) {
	spritesLen := len(h.Sprites)
	merged, mergedAlpha := make([]image.Image, spritesLen), make([]image.Image, spritesLen)

	base, err := decodeAsset(&h.Sprites[spritesLen-1], resDir, prefix)
	if err != nil {
		return nil, nil, err
	}
	baseAlpha, err := decodeAsset(&h.SpritesAlpha[spritesLen-1], resDir, prefix)
	if err != nil {
		return nil, nil, err
	}

	for i := range h.Sprites {
		// If i is the last index of h.Sprites, it is the base index, so we don't merge it and just use base instead.
		if i == spritesLen-1 {
			merged[i], mergedAlpha[i] = base, baseAlpha
			continue
		}

		face, err := decodeAsset(&h.Sprites[i], resDir, prefix)
		if err != nil {
			return nil, nil, err
		}
		faceAlpha, err := decodeAsset(&h.SpritesAlpha[i], resDir, prefix)
		if err != nil {
			return nil, nil, err
		}

		// If i is whole body, we don't merge it and just use face instead.
		if h.SpritesIsWholeBody[i] {
			merged[i], mergedAlpha[i] = face, baseAlpha
			continue
		}

		merged[i] = MergeFace(base, face, h.FacePos, h.FaceSize)
		mergedAlpha[i] = MergeFace(baseAlpha, faceAlpha, h.FacePos, h.FaceSize)
	}
	return merged, mergedAlpha, nil
}
func MergeFace(base image.Image, face image.Image, facePos pos, faceSize pos) image.Image {
	if facePos == (pos{-1, -1}) {
		return face
	}

	canvas := image.NewRGBA(base.Bounds())
	faceBounds := image.Rectangle{
		Min: image.Point{},
		Max: image.Point{X: int(faceSize.X), Y: int(faceSize.Y)},
	}.Add(image.Point{
		X: int(facePos.X),
		Y: int(facePos.Y),
	})
	face = imaging.Resize(face, int(faceSize.X), int(faceSize.Y), imaging.Lanczos)
	draw.Draw(canvas, base.Bounds(), base, image.Point{}, draw.Over)
	draw.Draw(canvas, faceBounds, face, image.Point{}, draw.Over)

	return canvas
}

func MergeAlpha(c image.Image, a image.Image) image.Image {
	a = imaging.Resize(a, c.Bounds().Dx(), c.Bounds().Dy(), imaging.Lanczos)
	canvas := image.NewNRGBA(c.Bounds())
	for x := 0; x < canvas.Bounds().Dx(); x++ {
		for y := 0; y < canvas.Bounds().Dy(); y++ {
			rc, gc, bc, _ := c.At(x, y).RGBA()
			ra, ga, ba, _ := a.At(x, y).RGBA()
			canvas.SetNRGBA(x, y, color.NRGBA{
				R: uint8(rc >> 8),
				G: uint8(gc >> 8),
				B: uint8(bc >> 8),
				A: uint8(((ra + ga + ba) / 3) >> 8),
			})
		}
	}
	return canvas
}

func decodeAsset(a *Asset, resDir, prefix string) (image.Image, error) {
	img, err := decodeFile(a.FilePath(resDir, prefix))
	if err != nil {
		return nil, err
	}
	return img, nil
}
func decodeFile(filePath string) (image.Image, error) {
	af, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func() { _ = af.Close() }()
	img, _, err := image.Decode(af)
	if err != nil {
		return nil, err
	}
	return img, nil
}
