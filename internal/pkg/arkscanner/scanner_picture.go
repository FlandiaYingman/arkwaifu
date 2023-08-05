package arkscanner

import (
	"fmt"
	"github.com/flandiayingman/arkwaifu/internal/pkg/util/pathutil"
	"github.com/pkg/errors"
	"path"
	"path/filepath"
)

type PictureArt struct {
	ID   string
	Name string
	Kind string
}

func (a *PictureArt) Path() string {
	subdirectory := subdirectoryOfCategory(a.Kind)
	return path.Join("assets/torappu/dynamicassets", "avg", subdirectory, a.Name)
}

func subdirectoryOfCategory(category string) string {
	switch category {
	case "image":
		return "images"
	case "background":
		return "backgrounds"
	case "item":
		return "items"
	default:
		panic(fmt.Sprintf("unrecognized category %s", category))
	}
}

const PictureArtPath = "assets/torappu/dynamicassets/avg"

func (scanner *Scanner) ScanForPictureArts() ([]*PictureArt, error) {
	baseDir := filepath.Join(scanner.Root, PictureArtPath)
	imageArts, err := scanner.scanPictures(baseDir, "images", "image")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	backgroundArts, err := scanner.scanPictures(baseDir, "backgrounds", "background")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	itemArts, err := scanner.scanPictures(baseDir, "items", "item")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var arts []*PictureArt
	arts = append(arts, imageArts...)
	arts = append(arts, backgroundArts...)
	arts = append(arts, itemArts...)

	return arts, nil
}

func (scanner *Scanner) scanPictures(base string, sub string, kind string) ([]*PictureArt, error) {
	files, err := filepath.Glob(filepath.Join(base, sub, "*.png"))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var pictures []*PictureArt
	for _, file := range files {
		pictures = append(pictures, &PictureArt{
			ID:   pathutil.Stem(file),
			Name: filepath.Base(file),
			Kind: kind,
		})
	}

	return pictures, nil
}
