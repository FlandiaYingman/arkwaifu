package asset

import (
	"fmt"
	"path/filepath"

	"github.com/uptrace/bun"
)

type mAsset struct {
	bun.BaseModel `bun:"table:asset_assets,alias:aa"`

	Kind string `bun:"kind,pk"`
	Name string `bun:"name,pk"`

	KindSortID int    `bun:"kind_sort_id,type:integer"`
	NameSortID []byte `bun:"name_sort_id,type:bytea"`

	Variants []*mVariant `bun:"rel:has-many,join:kind=asset_kind,join:name=asset_name"`
}

func (m mAsset) String() string {
	return fmt.Sprintf("%s/%s", m.Kind, m.Name)
}

type mVariant struct {
	bun.BaseModel `bun:"table:asset_variants"`

	AssetKind string `bun:"asset_kind,pk"`
	AssetName string `bun:"asset_name,pk"`
	Variant   string `bun:"variant,pk"`
	Filename  string `bun:"filename"`

	KindSortID    int    `bun:"kind_sort_id,type:integer"`
	NameSortID    []byte `bun:"name_sort_id,type:bytea"`
	VariantSortID int    `bun:"variant_sort_id,type:integer"`
}

func (v mVariant) String() string {
	return fmt.Sprintf("%s/%s/%s", v.AssetKind, v.AssetName, v.Variant)
}

// Path returns the relative path to the variant file.
//
// The path of assets is "{variant}/{asset.kind}/{filename}".
// When the program want to find an asset file in a certain directory, it will check the path relative to the directory.
func (v mVariant) Path() string {
	return fmt.Sprintf("%s/%s/%s", v.Variant, v.AssetKind, v.Filename)
}

// FilePath returns the absolute path to the variant file.
//
// This is a shortcut for filepath.Join(dir, v.Path()).
func (v mVariant) FilePath(dirPath string) string {
	return filepath.Join(dirPath, v.Path())
}

type mKindName struct {
	bun.BaseModel `bun:"table:asset_kind_names"`

	KindName string `bun:"kind_name,pk"`
	SortID   int    `bun:"sort_id,autoincrement"`
}
type mVariantName struct {
	bun.BaseModel `bun:"table:asset_variant_names"`

	VariantName string `bun:"variant_name,pk"`
	SortID      int    `bun:"sort_id,autoincrement"`
}

func (m mKindName) String() string {
	return fmt.Sprintf("%s:%d", m.KindName, m.SortID)
}
func (m mVariantName) String() string {
	return fmt.Sprintf("%s:%d", m.VariantName, m.SortID)
}
