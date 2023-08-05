package arkprocessor

import (
	_ "github.com/chai2010/webp"
	"github.com/flandiayingman/arkwaifu/internal/pkg/arkscanner"
	"github.com/pkg/errors"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
)

type PictureArt arkscanner.PictureArt

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
		return nil, errors.Wrapf(err, "process picture art %s", art.ID)
	} else {
		return &PictureArtImage{
			Image: img,
			Art:   art,
		}, nil
	}
}

func (a *PictureArt) decode(root string) (image.Image, error) {
	artPath := filepath.Join(root, (*arkscanner.PictureArt)(a).Path())

	artFile, err := os.Open(artPath)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer func() { _ = artFile.Close() }()

	img, _, err := image.Decode(artFile)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return img, nil
}
