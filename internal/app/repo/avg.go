package repo

import (
	"arkwaifu/internal/app/entity"
	"context"
	"database/sql"
	"fmt"
	"github.com/uptrace/bun"
)

type AvgRepo struct {
	// conn is for injection. e.g., Atomic.
	conn *bun.DB
	// db is for normal database operations.
	db bun.IDB
}

func (r *AvgRepo) Atomic(ctx context.Context, atomicBlock func(r *AvgRepo) error) (err error) {
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

	newRepo := &AvgRepo{
		conn: r.conn,
		db:   tx,
	}
	err = atomicBlock(newRepo)
	return
}

func NewAvgRepo(db *bun.DB) (*AvgRepo, error) {
	_, err := db.NewCreateTable().
		Model((*entity.Avg)(nil)).
		IfNotExists().
		Exec(context.Background())
	if err != nil {
		return nil, err
	}
	return &AvgRepo{
		conn: db,
		db:   db,
	}, nil
}

func (r *AvgRepo) GetAvgs(ctx context.Context) ([]*entity.Avg, error) {
	var items []*entity.Avg
	err := r.db.
		NewSelect().
		Model(items).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *AvgRepo) GetAvgByID(ctx context.Context, id string) (*entity.Avg, error) {
	var item entity.Avg
	err := r.db.
		NewSelect().
		Model(&item).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *AvgRepo) UpsertAvg(ctx context.Context, avg entity.Avg) error {
	_, err := r.db.
		NewInsert().
		Model(&avg).
		On("CONFLICT (story_id) DO UPDATE").
		Set("story_id = EXCLUDED.story_id").
		Set("story_code = EXCLUDED.story_code").
		Set("story_name = EXCLUDED.story_name").
		Set("story_txt = EXCLUDED.story_txt").
		Set("avg_tag = EXCLUDED.avg_tag").
		Set("group_id = EXCLUDED.group_id").
		Exec(ctx)
	return err
}
