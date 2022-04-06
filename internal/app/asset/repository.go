package asset

import (
	"context"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/uptrace/bun"
)

type repo struct {
	bun.IDB
	DB *bun.DB
}

func NewRepo(db *bun.DB) (*repo, error) {
	r := repo{IDB: db, DB: db}
	err := r.DB.RunInTx(context.Background(), nil, func(ctx context.Context, tx bun.Tx) error {
		var err error
		_, err = db.NewCreateTable().
			Model((*modelKindName)(nil)).
			IfNotExists().
			Exec(context.Background())
		if err != nil {
			return err
		}
		_, err = db.NewCreateTable().
			Model((*modelVariantName)(nil)).
			IfNotExists().
			Exec(context.Background())
		if err != nil {
			return err
		}
		_, err = db.NewCreateTable().
			Model((*modelAsset)(nil)).
			IfNotExists().
			ForeignKey("(kind) REFERENCES asset_kind_names (kind_name) ON DELETE CASCADE").
			Exec(context.Background())
		if err != nil {
			return err
		}
		_, err = db.NewCreateTable().
			Model((*modelVariant)(nil)).
			IfNotExists().
			ForeignKey("(variant) REFERENCES asset_variant_names (variant_name) ON DELETE CASCADE").
			ForeignKey("(asset_kind, asset_name) REFERENCES asset_assets (kind, name) ").
			Exec(context.Background())
		if err != nil {
			return err
		}
		return nil
	})
	return &r, err
}

type modelAsset struct {
	bun.BaseModel `bun:"table:asset_assets"`
	Kind          string          `bun:"kind,pk"`
	Name          string          `bun:"name,pk"`
	Variants      []*modelVariant `bun:"rel:has-many,join:kind=asset_kind,join:name=asset_name"`
}
type modelVariant struct {
	bun.BaseModel `bun:"table:asset_variants"`
	AssetKind     string `bun:"asset_kind,pk"`
	AssetName     string `bun:"asset_name,pk"`
	Variant       string `bun:"variant,pk"`
	Filename      string `bun:"filename"`
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

func (r *repo) SelectAssets(ctx context.Context, kind *string) ([]modelAsset, error) {
	models := new([]modelAsset)
	query := r.
		NewSelect().
		Model(models).
		Relation("Variants", SortVariant).
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
		Relation("Variants", SortVariant).
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
