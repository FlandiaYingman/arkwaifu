package avg

import (
	"context"
	"github.com/uptrace/bun"
)

func sortAvg(query *bun.SelectQuery) *bun.SelectQuery {
	return query.Order("sort_id ASC")
}
func sortAsset(query *bun.SelectQuery) *bun.SelectQuery {
	return query.Order("pk ASC")
}

type Repo struct {
	bun.IDB
	DB *bun.DB
}

func NewRepo(db *bun.DB) (*Repo, error) {
	r := Repo{DB: db, IDB: db}
	err := r.DB.RunInTx(context.Background(), nil, func(ctx context.Context, tx bun.Tx) error {
		var err error
		_, err = db.NewCreateTable().
			Model((*ResVersion)(nil)).
			IfNotExists().
			ColumnExpr("CHECK (id = ?)", true).
			Exec(context.Background())
		if err != nil {
			return err
		}
		_, err = db.NewCreateTable().
			Model((*groupModel)(nil)).
			IfNotExists().
			Exec(context.Background())
		if err != nil {
			return err
		}
		_, err = db.NewCreateTable().
			Model((*storyModel)(nil)).
			IfNotExists().
			Exec(context.Background())
		if err != nil {
			return err
		}
		_, err = db.NewCreateTable().
			Model((*assetModel)(nil)).
			IfNotExists().
			Exec(context.Background())
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (r Repo) UpdateAvg(ctx context.Context, version string, gms []groupModel, sms []storyModel) error {
	return r.DB.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		r.IDB = tx
		var err error

		err = r.Truncate(ctx)
		if err != nil {
			return err
		}
		err = r.Truncate(ctx)
		if err != nil {
			return err
		}
		err = r.UpsertVersion(ctx, version)
		if err != nil {
			return err
		}

		err = r.InsertGroups(ctx, gms)
		if err != nil {
			return err
		}

		err = r.InsertStories(ctx, sms)
		if err != nil {
			return err
		}

		return nil
	})
}
