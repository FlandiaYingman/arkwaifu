package avg

import (
	"context"
	"github.com/flandiayingman/arkwaifu/internal/app/infra"
	"github.com/uptrace/bun"
)

type VersionRepo struct {
	infra.Repo
}

type ResVersion struct {
	bun.BaseModel `bun:"table:version"`

	// ID can only be true, because this table should only have one row.
	ID         bool   `bun:",pk"`
	ResVersion string `bun:""`
}

func NewVersionRepo(db *bun.DB) (*VersionRepo, error) {
	_, err := db.NewCreateTable().
		Model((*ResVersion)(nil)).
		IfNotExists().
		Exec(context.Background())
	if err != nil {
		return nil, err
	}
	return &VersionRepo{
		Repo: infra.NewRepo(db),
	}, nil
}

func (r *VersionRepo) GetVersion(ctx context.Context) (string, error) {
	var resVersion ResVersion
	exists, err := r.DB().
		NewSelect().
		Model(&resVersion).
		Exists(ctx)
	if err != nil {
		return "", err
	}
	if exists {
		err := r.DB().
			NewSelect().
			Model(&resVersion).
			Scan(ctx)
		if err != nil {
			return "", err
		}
		return resVersion.ResVersion, nil
	}
	return "", nil
}

func (r *VersionRepo) UpsertVersion(ctx context.Context, resVersion string) error {
	resVersionEntity := ResVersion{
		ID:         true,
		ResVersion: resVersion,
	}
	_, err := r.DB().
		NewInsert().
		Model(&resVersionEntity).
		On("CONFLICT (id) DO UPDATE").
		Set("res_version = EXCLUDED.res_version").
		Exec(ctx)
	return err
}
