package asset

import (
	"context"
	"github.com/flandiayingman/arkwaifu/internal/app/config"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"io"
	"os"
	"path/filepath"
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

type Asset struct {
	Kind     string    `json:"kind"`
	Name     string    `json:"name"`
	Variants []Variant `json:"variant"`
}
type Variant struct {
	Variant  string `json:"variant"`
	Filename string `json:"filename"`
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

	assetFilePath := filepath.Join(s.staticDir, variant, kind, v.Filename)
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
	v := fromVariantModel(*mv)
	return &v, nil
}

func (s *Service) PostVariant(ctx context.Context, kind, name string, variant Variant, file io.Reader) error {
	vm := modelVariant{
		AssetKind: kind,
		AssetName: name,
		Variant:   variant.Variant,
		Filename:  variant.Filename,
	}

	a, err := s.GetAsset(ctx, kind, name)
	if err != nil {
		return err
	}
	if a == nil {
		return errors.Errorf("asset %s/%s not found", kind, name)
	}

	dstPath := filepath.Join(s.staticDir, vm.Variant, vm.AssetKind, variant.Filename)
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return errors.Wrapf(err, "failed to create dst %s", dstPath)
	}

	_, err = io.Copy(dstFile, file)
	if err != nil {
		return errors.Wrapf(err, "failed to copy file to %s", dstPath)
	}

	err = s.repo.InsertVariant(ctx, vm)
	if err != nil {
		return errors.Wrap(err, "failed to insert variant")
	}

	return nil
}

func fromAssetModel(model modelAsset) Asset {
	vms := lo.Map(model.Variants, func(vmPtr *modelVariant, _ int) Variant {
		return fromVariantModel(*vmPtr)
	})
	return Asset{
		Kind:     model.Kind,
		Name:     model.Name,
		Variants: vms,
	}
}
func fromVariantModel(model modelVariant) Variant {
	return Variant{
		Variant:  model.Variant,
		Filename: model.Filename,
	}
}

func wrapIter[T any, R any](f func(T) R) func(T, int) R {
	return func(t T, _ int) R {
		return f(t)
	}
}
