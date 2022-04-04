package asset

import (
	"context"
	"github.com/pkg/errors"
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
			Model((*modelAsset)(nil)).
			IfNotExists().
			Exec(context.Background())
		if err != nil {
			return err
		}
		_, err = db.NewCreateTable().
			Model((*modelVariant)(nil)).
			IfNotExists().
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

func (r *repo) Truncate(ctx context.Context) error {
	return r.DB.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		var err error
		_, err = tx.
			NewTruncateTable().
			Model((*modelAsset)(nil)).
			Exec(ctx)
		if err != nil {
			return err
		}
		_, err = tx.
			NewTruncateTable().
			Model((*modelVariant)(nil)).
			Exec(ctx)
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *repo) InsertAsset(ctx context.Context, models ...modelAsset) error {
	_, err := r.
		NewInsert().
		Model(&models).
		Exec(ctx)
	return err
}
func (r *repo) InsertVariant(ctx context.Context, models ...modelVariant) error {
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
		Relation("Variants")
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
		Relation("Variants").
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
