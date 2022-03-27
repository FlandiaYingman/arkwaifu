package asset

import (
	"context"
	"github.com/flandiayingman/arkwaifu/internal/app/avg"
	"github.com/flandiayingman/arkwaifu/internal/app/config"
	"github.com/pkg/errors"
	. "github.com/szmcdull/glinq/unsafe"
	"path/filepath"
)

type Service struct {
	staticDir   string
	repo        *repo
	versionRepo *avg.VersionRepo
}

func NewService(conf *config.Config, assetRepo *repo, versionRepo *avg.VersionRepo) *Service {
	return &Service{
		staticDir:   filepath.Join(conf.ResourceLocation, "static"),
		repo:        assetRepo,
		versionRepo: versionRepo,
	}
}

type Asset struct {
	Kind     string `json:"kind"`
	Name     string `json:"name"`
	Variant  string `json:"variant"`
	FileName string `json:"fileName"`
}

func (s *Service) SetAssets(ctx context.Context, assets []Asset) error {
	models := ToSlice(Select(FromSlice(assets), toModel))

	err := s.repo.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = s.repo.EndTx(err) }()

	err = s.repo.Truncate(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to truncate asset table")
	}

	err = s.repo.Insert(ctx, models...)
	if err != nil {
		return errors.Wrap(err, "failed to insert assets")
	}

	return nil
}
func (s *Service) GetAssets(ctx context.Context, kind, name, variant string) ([]Asset, error) {
	kindPtr := &kind
	namePtr := &name
	variantPtr := &variant
	if *kindPtr == "" {
		kindPtr = nil
	}
	if *namePtr == "" {
		namePtr = nil
	}
	if *variantPtr == "" {
		variantPtr = nil
	}

	models, err := s.repo.SelectAll(ctx, kindPtr, namePtr, variantPtr)
	if err != nil {
		return nil, err
	}

	return ToSlice(Select(FromSlice(models), fromModel)), nil
}
func (s *Service) GetAsset(ctx context.Context, kind, name, variant string) (Asset, error) {
	kindPtr := &kind
	namePtr := &name
	variantPtr := &variant
	if *kindPtr == "" {
		kindPtr = nil
	}
	if *namePtr == "" {
		namePtr = nil
	}
	if *variantPtr == "" {
		variantPtr = nil
	}

	m, err := s.repo.SelectOne(ctx, kindPtr, namePtr, variantPtr)
	if err != nil {
		return Asset{}, err
	}

	return fromModel(m), nil
}
func (s *Service) GetKinds(ctx context.Context) ([]string, error) {
	kinds, err := s.repo.SelectUniqueKinds(ctx, nil, nil, nil)
	return kinds, err
}
func (s *Service) GetNames(ctx context.Context, kind string) ([]string, error) {
	names, err := s.repo.SelectUniqueNames(ctx, &kind, nil, nil)
	return names, err
}
func (s *Service) GetVariants(ctx context.Context, kind, name string) ([]string, error) {
	variants, err := s.repo.SelectUniqueVariants(ctx, &kind, &name, nil)
	return variants, err
}
func (s *Service) GetAssetFilePath(ctx context.Context, kind, name, variant string) (string, error) {
	asset, err := s.GetAsset(ctx, kind, name, variant)
	if err != nil {
		return "", err
	}

	assetFilePath := filepath.Join(s.staticDir, asset.Variant, asset.Kind, asset.FileName)
	return assetFilePath, nil
}

func toModel(asset Asset) model {
	return model{
		Kind:     asset.Kind,
		Name:     asset.Name,
		Variant:  asset.Variant,
		FileName: asset.FileName,
	}
}
func fromModel(model model) Asset {
	return Asset{
		Kind:     model.Kind,
		Name:     model.Name,
		Variant:  model.Variant,
		FileName: model.FileName,
	}
}
