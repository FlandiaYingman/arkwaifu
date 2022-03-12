package res

import (
	"context"
	"errors"
	"github.com/flandiayingman/arkwaifu/internal/app/avg"
	"github.com/flandiayingman/arkwaifu/internal/app/config"
	"os"
	"path"
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

func (s *Service) GetImages(ctx context.Context) ([]Resource, error) {
	location, err := s.getImageLocation(ctx, Raw)
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(location)
	if err != nil {
		return nil, err
	}
	resources := make([]Resource, len(entries))
	for i, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return nil, nil
			}
			return nil, err
		}
		resources[i] = *getResourceByFileInfo(location, info)
	}
	return resources, nil
}

func (s *Service) GetImageByName(ctx context.Context, name string, variant Variant) (*Resource, error) {
	location, err := s.getImageLocation(ctx, variant)
	if err != nil {
		return nil, err
	}
	filename := variant.FileName(name)
	info, err := os.Stat(filepath.Join(location, filename))
	if err == nil {
		return getResourceByFileInfo(location, info), nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	return nil, err
}

func (s *Service) getImageLocation(ctx context.Context, variant Variant) (string, error) {
	version, err := s.versionRepo.GetVersion(ctx)
	if err != nil {
		return "", err
	}
	location := filepath.Join(variant.Location(s.resourceLocation, version), "images")
	return location, nil
}

func (s *Service) GetBackgrounds(ctx context.Context) ([]Resource, error) {
	location, err := s.getBackgroundLocation(ctx, Raw)
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(location)
	if err != nil {
		return nil, err
	}
	resources := make([]Resource, len(entries))
	for i, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return nil, nil
			}
			return nil, err
		}
		resources[i] = *getResourceByFileInfo(location, info)
	}
	return resources, nil
}

func (s *Service) GetBackgroundByName(ctx context.Context, name string, variant Variant) (*Resource, error) {
	location, err := s.getBackgroundLocation(ctx, variant)
	if err != nil {
		return nil, err
	}
	filename := variant.FileName(name)
	info, err := os.Stat(filepath.Join(location, filename))
	if err == nil {
		return getResourceByFileInfo(location, info), nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	return nil, err
}

func (s *Service) getBackgroundLocation(ctx context.Context, variant Variant) (string, error) {
	version, err := s.versionRepo.GetVersion(ctx)
	if err != nil {
		return "", err
	}
	location := filepath.Join(variant.Location(s.resourceLocation, version), "backgrounds")
	return location, nil
}

func getResourceByFileInfo(location string, info os.FileInfo) *Resource {
	filename := info.Name()
	resource := Resource{
		Name: strings.TrimSuffix(filename, path.Ext(filename)),
		Path: filepath.Join(location, filename),
	}
	return &resource
}
