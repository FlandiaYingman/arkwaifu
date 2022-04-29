package asset

import (
	"fmt"
	"path/filepath"

	"github.com/samber/lo"
)

const (
	KindImages      = "images"
	KindBackgrounds = "backgrounds"
	KindCharacters  = "characters"

	VariantImg        = "img"
	VariantTimg       = "timg"
	VariantAlpha      = "alpha"
	VariantRealEsrgan = "real-esrgan"
	VariantRealCugan  = "real-cugan"
)

var (
	Kinds = []string{
		KindImages,
		KindBackgrounds,
		KindCharacters,
	}
	KindsMap = map[string]string{
		KindImages:      KindImages,
		KindBackgrounds: KindBackgrounds,
		KindCharacters:  KindCharacters,
	}
	Variants = []string{
		VariantImg,
		VariantTimg,
		VariantAlpha,
		VariantRealEsrgan,
		VariantRealCugan,
	}
	VariantsMap = map[string]string{
		VariantImg:        VariantImg,
		VariantTimg:       VariantTimg,
		VariantAlpha:      VariantAlpha,
		VariantRealEsrgan: VariantRealEsrgan,
		VariantRealCugan:  VariantRealCugan,
	}
)

var (
	ErrInvalidKind    = fmt.Errorf("invalid kind")
	ErrInvalidVariant = fmt.Errorf("invalid variant")
)

func CheckKind(kind string) error {
	if _, ok := KindsMap[kind]; !ok {
		return fmt.Errorf("%s: %w", kind, ErrInvalidKind)
	}
	return nil
}
func CheckVariant(variant string) error {
	if _, ok := VariantsMap[variant]; !ok {
		return fmt.Errorf("%s: %w", variant, ErrInvalidVariant)
	}
	return nil
}

type Asset struct {
	Kind     string     `json:"kind"`
	Name     string     `json:"name"`
	Variants *[]Variant `json:"variants"`
}
type Variant struct {
	Variant  string `json:"variant"`
	Filename string `json:"filename"`
	Asset    *Asset `json:"-"`
}

func fromAssetModel(model modelAsset) Asset {
	asset := Asset{
		Kind:     model.Kind,
		Name:     model.Name,
		Variants: nil,
	}
	vs := lo.Map(model.Variants,
		func(vmPtr *modelVariant, _ int) Variant {
			v := fromVariantModel(*vmPtr)
			v.Asset = &asset
			return v
		},
	)
	asset.Variants = &vs
	return asset
}
func fromVariantModel(model modelVariant) Variant {
	return Variant{
		Variant:  model.Variant,
		Filename: model.Filename,
		Asset:    nil,
	}
}

// String returns a string representation of the asset.
//
// The string representation of assets is "{kind}/{name}".
// e.g., "images/20_i00", "backgrounds/21_g1_interrogat_room".
func (a Asset) String() string {
	return fmt.Sprintf("%s/%s", a.Kind, a.Name)
}

// String returns a string representation of the variant.
//
// The string representation of assets is "{asset}/{variant}".
// e.g., "images/20_i00/img", "backgrounds/21_g1_interrogat_room/timg".
func (v Variant) String() string {
	return fmt.Sprintf("%s/%s", v.Asset.String(), v.Variant)
}

// Path returns the relative path to the variant file.
//
// The path of assets is "{variant}/{asset.kind}/{filename}".
// When the program want to find an asset file in a certain directory, it will check the path relative to the directory.
func (v Variant) Path() string {
	return fmt.Sprintf("%s/%s/%s", v.Variant, v.Asset.Kind, v.Filename)
}

// FilePath returns the absolute path to the variant file.
//
// This is a shortcut for filepath.Join(dir, v.Path()).
func (v Variant) FilePath(dirPath string) string {
	return filepath.Join(dirPath, v.Path())
}
