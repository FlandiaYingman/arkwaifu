package asset

import (
	"context"
	"github.com/flandiayingman/arkwaifu/internal/app/infra"
	"github.com/samber/lo"
	"github.com/uptrace/bun"
)

type model struct {
	bun.BaseModel `bun:"table:asset_assets"`
	Kind          string `bun:"kind,pk"`
	Name          string `bun:"name,pk"`
	Variant       string `bun:"variant,pk"`
	FileName      string `bun:"fileName"`
}

type repo struct {
	infra.Repo
}

func NewRepo(db *bun.DB) (*repo, error) {
	_, err := db.NewCreateTable().
		Model((*model)(nil)).
		IfNotExists().
		Exec(context.Background())
	if err != nil {
		return nil, err
	}
	r := repo{Repo: infra.NewRepo(db)}
	return &r, nil
}

func (r *repo) Truncate(ctx context.Context) error {
	_, err := r.DB().
		NewTruncateTable().
		Model((*model)(nil)).
		Exec(ctx)
	return err
}
func (r *repo) Insert(ctx context.Context, models ...model) error {
	_, err := r.DB().
		NewInsert().
		Model(&models).
		Exec(ctx)
	return err
}
func (r *repo) SelectAll(ctx context.Context, kind, asset, variant *string) ([]model, error) {
	var models []model
	query := r.DB().
		NewSelect().
		Model(&models)
	if kind != nil {
		query = query.Where("kind = ?", *kind)
	}
	if asset != nil {
		query = query.Where("name = ?", *asset)
	}
	if variant != nil {
		query = query.Where("variant = ?", *variant)
	}
	err := query.Scan(ctx)
	return models, err
}
func (r *repo) SelectOne(ctx context.Context, kind, asset, variant *string) (model, error) {
	var models model
	query := r.DB().
		NewSelect().
		Model(&models)
	if kind != nil {
		query = query.Where("kind = ?", *kind)
	}
	if asset != nil {
		query = query.Where("name = ?", *asset)
	}
	if variant != nil {
		query = query.Where("variant = ?", *variant)
	}
	err := query.Scan(ctx)
	return models, err
}
func (r *repo) SelectUniqueKinds(ctx context.Context, kind, asset, variant *string) ([]string, error) {
	var models []model
	query := r.DB().
		NewSelect().
		DistinctOn("(kind)").
		Model(&models)
	if kind != nil {
		query = query.Where("kind = ?", *kind)
	}
	if asset != nil {
		query = query.Where("name = ?", *asset)
	}
	if variant != nil {
		query = query.Where("variant = ?", *variant)
	}
	err := query.Scan(ctx)
	return lo.Map(models, func(m model, _ int) string { return m.Kind }), err
}
func (r *repo) SelectUniqueNames(ctx context.Context, kind, asset, variant *string) ([]string, error) {
	var models []model
	query := r.DB().
		NewSelect().
		DistinctOn("(name)").
		Model(&models)
	if kind != nil {
		query = query.Where("kind = ?", *kind)
	}
	if asset != nil {
		query = query.Where("name = ?", *asset)
	}
	if variant != nil {
		query = query.Where("variant = ?", *variant)
	}
	err := query.Scan(ctx)
	return lo.Map(models, func(m model, _ int) string { return m.Name }), err
}
func (r *repo) SelectUniqueVariants(ctx context.Context, kind, asset, variant *string) ([]string, error) {
	var models []model
	query := r.DB().
		NewSelect().
		DistinctOn("(variant)").
		Model(&models)
	if kind != nil {
		query = query.Where("kind = ?", *kind)
	}
	if asset != nil {
		query = query.Where("name = ?", *asset)
	}
	if variant != nil {
		query = query.Where("variant = ?", *variant)
	}
	err := query.Scan(ctx)
	return lo.Map(models, func(m model, _ int) string { return m.Variant }), err
}
