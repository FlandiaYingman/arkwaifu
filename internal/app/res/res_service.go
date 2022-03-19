package res

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

type Resource struct {
	Name string
	Path string
}

func NewService(conf *config.Config, repo *avg.VersionRepo) *Service {
	return &Service{
		resourceLocation: conf.ResourceLocation,
		versionRepo:      repo,
	}
}

func (s *Service) GetAssets(variant Variant, kind Kind) ([]string, error) {
	dirPath := filepath.Join(s.resourceLocation, "static", string(variant), string(kind))
	dir, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	var result []string
	linq.From(dir).Select(func(i interface{}) interface{} {
		name := i.(os.DirEntry).Name()
		return strings.TrimSuffix(name, filepath.Ext(name))
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
