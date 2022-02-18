package repo

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/uptrace/bun"
)

type ResVersionRepo struct {
	// conn is for injection. e.g., Atomic.
	conn *bun.DB
	// db is for normal database operations.
	db bun.IDB
}

type ResVersion struct {
	bun.BaseModel `bun:"table:res_version"`

	// ID can only be true, because this table should only have one row.
	ID         bool   `bun:",pk"`
	ResVersion string `bun:""`
}

func (r *ResVersionRepo) Atomic(ctx context.Context, atomicBlock func(r *ResVersionRepo) error) (err error) {
	tx, err := r.conn.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				err = fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
			}
		} else {
			err = tx.Commit()
		}
	}()

	newRepo := &ResVersionRepo{
		conn: r.conn,
		db:   tx,
	}
	err = atomicBlock(newRepo)
	return
}

func NewResVersionRepo(db *bun.DB) (*ResVersionRepo, error) {
	_, err := db.NewCreateTable().
		Model((*ResVersion)(nil)).
		IfNotExists().
		Exec(context.Background())
	if err != nil {
		return nil, err
	}
	return &ResVersionRepo{
		conn: db,
		db:   db,
	}, nil
}

func (r *ResVersionRepo) GetResVersion(ctx context.Context) (string, error) {
	var resVersion ResVersion
	exists, err := r.db.
		NewSelect().
		Model(&resVersion).
		Exists(ctx)
	if err != nil {
		return "", err
	}
	if exists {
		err := r.db.
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

func (r *ResVersionRepo) UpsertResVersion(ctx context.Context, resVersion string) error {
	resVersionEntity := ResVersion{
		ID:         true,
		ResVersion: resVersion,
	}
	_, err := r.db.
		NewInsert().
		Model(&resVersionEntity).
		On("CONFLICT (id) DO UPDATE").
		Set("res_version = EXCLUDED.res_version").
		Exec(ctx)
	return err
}
