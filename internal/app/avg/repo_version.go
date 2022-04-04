package avg

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

type ResVersion struct {
	bun.BaseModel `bun:"table:version"`

	// ID can only be true, because this table should only have one row.
	ID         *bool  `bun:",pk,default:true"`
	ResVersion string `bun:""`
}

func (r *Repo) GetVersion(ctx context.Context) (string, error) {
	var resVersion ResVersion
	err := r.
		NewSelect().
		Model(&resVersion).
		Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return resVersion.ResVersion, nil
}
func (r *Repo) UpsertVersion(ctx context.Context, resVersion string) error {
	resVersionEntity := ResVersion{
		ResVersion: resVersion,
	}
	_, err := r.
		NewInsert().
		Model(&resVersionEntity).
		On("CONFLICT (id) DO UPDATE").
		Set("res_version = EXCLUDED.res_version").
		Exec(ctx)
	return err
}
