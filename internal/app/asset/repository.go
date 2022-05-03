package asset

import (
	"context"
	"fmt"
	"io"
	"path/filepath"

	"github.com/flandiayingman/arkwaifu/internal/app/config"
	"github.com/flandiayingman/arkwaifu/internal/pkg/util/fileutil"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/uptrace/bun"
)

type repo struct {
	bun.IDB
	DB *bun.DB

	StaticDir string
}

func NewRepo(db *bun.DB, conf *config.Config) (*repo, error) {
	r := repo{
		IDB:       db,
		DB:        db,
		StaticDir: conf.StaticDir,
	}

	err := r.init()
	if err != nil {
		return nil, err
	}

	return &r, nil
}

type modelAsset struct {
	bun.BaseModel `bun:"table:asset_assets,alias:aa"`

	Kind string `bun:"kind,pk"`
	Name string `bun:"name,pk"`

	KindSortID int    `bun:"kind_sort_id,type:integer"`
	NameSortID []byte `bun:"name_sort_id,type:bytea"`

	Variants []*modelVariant `bun:"rel:has-many,join:kind=asset_kind,join:name=asset_name"`
}
type modelVariant struct {
	bun.BaseModel `bun:"table:asset_variants"`

	AssetKind string `bun:"asset_kind,pk"`
	AssetName string `bun:"asset_name,pk"`
	Variant   string `bun:"variant,pk"`
	Filename  string `bun:"filename"`

	KindSortID    int    `bun:"kind_sort_id,type:integer"`
	NameSortID    []byte `bun:"name_sort_id,type:bytea"`
	VariantSortID int    `bun:"variant_sort_id,type:integer"`
}
type modelKindName struct {
	bun.BaseModel `bun:"table:asset_kind_names"`
	KindName      string `bun:"kind_name,pk"`
	SortID        int    `bun:"sort_id,autoincrement"`
}
type modelVariantName struct {
	bun.BaseModel `bun:"table:asset_variant_names"`
	VariantName   string `bun:"variant_name,pk"`
	SortID        int    `bun:"sort_id,autoincrement"`
}

var (
	ErrExists = errors.New("asset: the asset or variant already exists.")
)

func (r *repo) Truncate(ctx context.Context) error {
	return r.DB.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		var err error
		_, err = tx.
			NewTruncateTable().
			Model((*modelAsset)(nil)).
			Cascade().
			Exec(ctx)
		if err != nil {
			return err
		}
		_, err = tx.
			NewTruncateTable().
			Model((*modelVariant)(nil)).
			Cascade().
			Exec(ctx)
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *repo) InsertAsset(ctx context.Context, models ...modelAsset) error {
	if len(models) == 0 {
		return nil
	}
	_, err := r.
		NewInsert().
		Model(&models).
		Exec(ctx)
	return err
}
func (r *repo) InsertVariant(ctx context.Context, models ...modelVariant) error {
	if len(models) == 0 {
		return nil
	}
	_, err := r.
		NewInsert().
		Model(&models).
		Exec(ctx)
	return err
}
func (r *repo) InsertVariantFile(ctx context.Context, m modelVariant, f io.Reader) error {
	return r.DB.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		// Insert the variant. The variant is invisible until the transaction is committed.
		res, err := tx.NewInsert().
			Model(&m).
			On("CONFLICT DO NOTHING").
			Exec(ctx)
		if err != nil {
			return err
		}

		// Check if the insertion conflicts. If it conflicts, return ErrExists.
		rows, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if rows <= 0 {
			return ErrExists
		}

		// Write variant file. Note MkFileFromReader closes the reader.
		err = fileutil.MkFileFromReader(m.FilePath(r.StaticDir), f)
		if err != nil {
			return err
		}

		// Return nil to commit if everything success.
		return nil
	})
}

func (r *repo) SelectAssets(ctx context.Context, kind *string) ([]modelAsset, error) {
	models := new([]modelAsset)
	query := r.
		NewSelect().
		Model(models).
		Relation("Variants", SortAssetVariant).
		Apply(SortAsset)
	if kind != nil {
		query.Where("kind = ?", *kind)
	}
	err := query.Scan(ctx)
	return *models, err
}
func (r *repo) SelectAsset(ctx context.Context, kind, name string) (*modelAsset, error) {
	model := new(modelAsset)
	err := r.
		NewSelect().
		Model(model).
		Relation("Variants", SortAssetVariant).
		Where("(kind, name) = (?, ?)", kind, name).
		Scan(ctx)
	return model, err
}

func (r *repo) SelectVariants(ctx context.Context, kind, name string) ([]modelVariant, error) {
	models := new([]modelVariant)
	err := r.
		NewSelect().
		Model(models).
		Where("(asset_kind, asset_name) = (?, ?)", kind, name).
		Apply(SortVariant).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return *models, err
}
func (r *repo) SelectVariant(ctx context.Context, kind, name, variant string) (*modelVariant, error) {
	model := new(modelVariant)
	err := r.
		NewSelect().
		Model(model).
		Where("(asset_kind, asset_name, variant) = (?, ?, ?)", kind, name, variant).
		Scan(ctx)
	return model, err
}

func (r *repo) InitNames(ctx context.Context, kindNames []string, variantNames []string) error {
	kms := lo.Map(kindNames, func(s string, _ int) modelKindName { return modelKindName{KindName: s} })
	vms := lo.Map(variantNames, func(s string, _ int) modelVariantName { return modelVariantName{VariantName: s} })
	return r.DB.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		r.IDB = tx
		var err error
		_, err = tx.
			NewTruncateTable().
			Model((*modelKindName)(nil)).
			Cascade().
			Exec(ctx)
		if err != nil {
			return err
		}
		_, err = tx.
			NewTruncateTable().
			Model((*modelVariantName)(nil)).
			Cascade().
			Exec(ctx)
		if err != nil {
			return err
		}
		_, err = r.
			NewInsert().
			Model(&kms).
			Exec(ctx)
		if err != nil {
			return err
		}
		_, err = r.
			NewInsert().
			Model(&vms).
			Exec(ctx)
		if err != nil {
			return err
		}
		return nil
	})
}
func (r *repo) SelectKindNames(ctx context.Context) ([]string, error) {
	models := new([]modelKindName)
	err := r.
		NewSelect().
		Model(models).
		Order("sort_id").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return lo.Map(*models, func(m modelKindName, _ int) string { return m.KindName }), nil
}
func (r *repo) SelectVariantNames(ctx context.Context) ([]string, error) {
	models := new([]modelVariantName)
	err := r.
		NewSelect().
		Model(models).
		Order("sort_id").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return lo.Map(*models, func(m modelVariantName, _ int) string { return m.VariantName }), nil
}

func (r repo) Update(ctx context.Context, ams []modelAsset, vms []modelVariant) error {
	return r.DB.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		r.IDB = tx
		var err error
		err = r.Truncate(ctx)
		if err != nil {
			return errors.Wrapf(err, "failed to truncate asset table")
		}
		err = r.InsertAsset(ctx, ams...)
		if err != nil {
			return errors.Wrapf(err, "failed to insert assets %v", ams)
		}
		err = r.InsertVariant(ctx, vms...)
		if err != nil {
			return errors.Wrapf(err, "failed to insert variants %v", vms)
		}
		return nil
	})
}

// Path returns the relative path to the variant file.
//
// The path of assets is "{variant}/{asset.kind}/{filename}".
// When the program want to find an asset file in a certain directory, it will check the path relative to the directory.
func (v modelVariant) Path() string {
	return fmt.Sprintf("%s/%s/%s", v.Variant, v.AssetKind, v.Filename)
}

// FilePath returns the absolute path to the variant file.
//
// This is a shortcut for filepath.Join(dir, v.Path()).
func (v modelVariant) FilePath(dirPath string) string {
	return filepath.Join(dirPath, v.Path())
}
