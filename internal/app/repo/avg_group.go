package repo

import (
	"arkwaifu/internal/app/entity"
	"context"
	"database/sql"
	"fmt"
	"github.com/uptrace/bun"
)

type AvgGroupRepo struct {
	// conn is for injection. e.g., Atomic.
	conn *bun.DB
	// db is for normal database operations.
	db bun.IDB
}

func (r *AvgGroupRepo) Atomic(ctx context.Context, atomicBlock func(r *AvgGroupRepo) error) (err error) {
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

	newRepo := &AvgGroupRepo{
		conn: r.conn,
		db:   tx,
	}
	err = atomicBlock(newRepo)
	return
}

func NewAvgGroupRepo(db *bun.DB) (*AvgGroupRepo, error) {
	_, err := db.NewCreateTable().
		Model((*entity.AvgGroup)(nil)).
		IfNotExists().
		Exec(context.Background())
	if err != nil {
		return nil, err
	}
	return &AvgGroupRepo{
		db:   db,
		conn: db,
	}, nil
}

func (r *AvgGroupRepo) GetAvgGroups(ctx context.Context) ([]*entity.AvgGroup, error) {
	var items []*entity.AvgGroup
	err := r.db.
		NewSelect().
		Model(items).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *AvgGroupRepo) GetAvgGroupByID(ctx context.Context, id string) (*entity.AvgGroup, error) {
	var item entity.AvgGroup
	err := r.db.
		NewSelect().
		Model(&item).
		Where("id = ?", id).
		Relation("Avgs").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *AvgGroupRepo) UpsertAvgGroup(ctx context.Context, group entity.AvgGroup) error {
	_, err := r.db.
		NewInsert().
		Model(&group).
		On("CONFLICT (id) DO UPDATE").
		Set("name = EXCLUDED.name").
		Exec(ctx)
	return err
}
