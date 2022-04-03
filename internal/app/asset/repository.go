package asset

import (
	"context"
	"github.com/flandiayingman/arkwaifu/internal/app/infra"
	"github.com/uptrace/bun"
)

type repo struct {
	infra.Repo
}

func NewRepo(db *bun.DB) (*repo, error) {
	var err error
	_, err = db.NewCreateTable().
		Model((*modelAsset)(nil)).
		IfNotExists().
		Exec(context.Background())
	if err != nil {
		return nil, err
	}
	_, err = db.NewCreateTable().
		Model((*modelVariant)(nil)).
		IfNotExists().
		Exec(context.Background())
	if err != nil {
		return nil, err
	}
	r := repo{Repo: infra.NewRepo(db)}
	return &r, nil
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
	var err error
	_, err = r.DB().
		NewTruncateTable().
		Model((*modelAsset)(nil)).
		Exec(ctx)
	if err != nil {
		return err
	}
	_, err = r.DB().
		NewTruncateTable().
		Model((*modelVariant)(nil)).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *repo) InsertAsset(ctx context.Context, models ...modelAsset) error {
	_, err := r.DB().
		NewInsert().
		Model(&models).
		Exec(ctx)
	return err
}
func (r *repo) InsertVariant(ctx context.Context, models ...modelVariant) error {
	_, err := r.DB().
		NewInsert().
		Model(&models).
		Exec(ctx)
	return err
}

func (r *repo) SelectAssets(ctx context.Context, kind *string) ([]modelAsset, error) {
	models := new([]modelAsset)
	query := r.DB().
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
	err := r.DB().
		NewSelect().
		Model(model).
		Relation("Variants").
		Where("(kind, name) = (?, ?)", kind, name).
		Scan(ctx)
	return model, err
}

func (r *repo) SelectVariants(ctx context.Context, kind, name string) ([]modelVariant, error) {
	models := new([]modelVariant)
	query := r.DB().
		NewSelect().
		Model(models).
		Where("(asset_kind, asset_name) = (?, ?)", kind, name).
		Scan(ctx)
	err := query
	if err != nil {
		return nil, err
	}
	return *models, err
}
func (r *repo) SelectVariant(ctx context.Context, kind, name, variant string) (*modelVariant, error) {
	model := new(modelVariant)
	err := r.DB().
		NewSelect().
		Model(model).
		Where("(asset_kind, asset_name, variant) = (?, ?, ?)", kind, name, variant).
		Scan(ctx)
	return model, err
}
