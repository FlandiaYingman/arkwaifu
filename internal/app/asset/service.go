package asset

import (
	"context"
	"io"

	"github.com/flandiayingman/arkwaifu/internal/app/config"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

type Service struct {
	staticDir string
	repo      *repo
}

func NewService(conf *config.Config, assetRepo *repo) *Service {
	return &Service{
		staticDir: conf.StaticDir,
		repo:      assetRepo,
	}
}

func (s *Service) GetAssets(ctx context.Context, kind *string) ([]Asset, error) {
	models, err := s.repo.SelectAssets(ctx, kind)
	if err != nil {
		return nil, err
	}
	return lo.Map(models, wrapIter(fromAssetModel)), nil
}
func (s *Service) GetAsset(ctx context.Context, kind, name string) (*Asset, error) {
	m, err := s.repo.SelectAsset(ctx, kind, name)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, nil
	}

	asset := fromAssetModel(*m)
	return &asset, nil
}

func (s *Service) GetVariantFile(ctx context.Context, kind, name, variant string) (*string, error) {
	v, err := s.GetVariant(ctx, kind, name, variant)
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, nil
	}

	assetFilePath := v.FilePath(s.staticDir)
	return &assetFilePath, nil
}
func (s *Service) GetVariants(ctx context.Context, kind string, name string) ([]Variant, error) {
	mVariants, err := s.repo.SelectVariants(ctx, kind, name)
	if err != nil {
		return nil, err
	}
	return lo.Map(mVariants, wrapIter(fromVariantModel)), nil
}
func (s *Service) GetVariant(ctx context.Context, kind string, name string, variant string) (*Variant, error) {
	mv, err := s.repo.SelectVariant(ctx, kind, name, variant)
	if err != nil {
		return nil, err
	}
	ma, err := s.repo.SelectAsset(ctx, kind, name)
	if err != nil {
		return nil, err
	}
	v := fromVariantModel(*mv)
	a := fromAssetModel(*ma)
	v.Asset = &a
	return &v, nil
}

func (s *Service) InitNames(ctx context.Context) error {
	return s.repo.InitNames(ctx, Kinds, Variants)
}
func (s *Service) GetKindNames(ctx context.Context) ([]string, error) {
	kinds, err := s.repo.SelectKindNames(ctx)
	if err != nil {
		return nil, err
	}
	return kinds, nil
}
func (s *Service) GetVariantNames(ctx context.Context) ([]string, error) {
	variants, err := s.repo.SelectVariantNames(ctx)
	if err != nil {
		return nil, err
	}
	return variants, nil
}

func (s *Service) PostVariant(ctx context.Context, kind, name string, variant Variant, file io.Reader) error {
	vm := mVariant{
		AssetKind: kind,
		AssetName: name,
		Variant:   variant.Variant,
		Filename:  variant.Filename,
	}
	err := s.repo.InsertVariantFile(ctx, vm, file)
	if err != nil {
		return errors.Wrap(err, "failed to insert variant")
	}
	return nil
}

func wrapIter[T any, R any](f func(T) R) func(T, int) R {
	return func(t T, _ int) R {
		return f(t)
	}
}
