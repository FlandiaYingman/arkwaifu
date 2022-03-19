package asset

import (
	"github.com/ahmetb/go-linq/v3"
	"github.com/flandiayingman/arkwaifu/internal/app/avg"
	"github.com/flandiayingman/arkwaifu/internal/app/config"
	"os"
	"path/filepath"
	"strings"
)

type Service struct {
	resourceLocation string
	versionRepo      *avg.VersionRepo
}

type Asset struct {
	ID      string  `json:"id"`
	Variant Variant `json:"variant"`
	Kind    Kind    `json:"kind"`
}

func NewService(conf *config.Config, repo *avg.VersionRepo) *Service {
	return &Service{
		resourceLocation: conf.ResourceLocation,
		versionRepo:      repo,
	}
}

func (s *Service) GetAssets(variant *Variant, kind *Kind) ([]Asset, error) {
	var assets []Asset
	if variant == nil {
		for _, v := range Variants {
			vAssets, err := s.GetAssets(&v, kind)
			if err != nil {
				return nil, err
			}

			assets = append(assets, vAssets...)
		}
		return assets, nil
	}
	if kind == nil {
		for _, k := range Kinds {
			kAssets, err := s.GetAssets(variant, &k)
			if err != nil {
				return nil, err
			}

			assets = append(assets, kAssets...)
		}
		return assets, nil
	}
	dirPath := filepath.Join(s.resourceLocation, "static", string(*variant), string(*kind))
	dir, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	var result []Asset
	linq.From(dir).Select(func(i interface{}) interface{} {
		name := i.(os.DirEntry).Name()
		name = strings.TrimSuffix(name, filepath.Ext(name))
		return Asset{
			ID:      name,
			Variant: *variant,
			Kind:    *kind,
		}
	}).ToSlice(&result)
	return result, nil
}

func (s *Service) GetAssetByID(id string, variant Variant, kind Kind) (string, error) {
	name := id + VariantExts[variant]
	path := filepath.Join(s.resourceLocation, "static", string(variant), string(kind), name)

	_, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	return path, nil
}
