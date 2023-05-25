package ark

import (
	"fmt"
	_ "github.com/chai2010/webp"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
)

type PictureArtImage struct {
	Image image.Image
	Art   *PictureArt
}

// ProcessPictureArt process the picture art.
//
// Since picture arts are trivial and different from character arts, the only
// thing this method does is to read the art image and return it.
func (p *Processor) ProcessPictureArt(art *PictureArt) (*PictureArtImage, error) {
	img, err := art.decode(p.Root)
	if err != nil {
		return nil, fmt.Errorf("process picture art %s: %w", art.ID, err)
	} else {
		return &PictureArtImage{
			Image: img,
			Art:   art,
		}, nil
	}
}

func (a *PictureArt) decode(root string) (image.Image, error) {
	artPath := filepath.Join(root, a.Path())

	artFile, err := os.Open(artPath)
	if err != nil {
		return nil, err
	}
	defer func() { _ = artFile.Close() }()

	img, _, err := image.Decode(artFile)
	if err != nil {
		return nil, err
	}

	return img, nil
}
