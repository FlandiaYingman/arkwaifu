package repo

import (
	"context"
	"database/sql"
	"github.com/uptrace/bun"
)

// Avg is a fragment of story. e.g.: "8-1 行动前" or "IW-9 行动后" (IW stands for activity "将进酒").
type Avg struct {
	bun.BaseModel `bun:"table:avgs"`

	StoryID   string `bun:",pk"`
	StoryCode string
	StoryName string
	StoryTxt  string
	AvgTag    string

	GroupID string `bun:"group_id"`
}

type AvgRepo struct {
	db *bun.DB
}

func NewAvgRepo(db *bun.DB) (*AvgRepo, error) {
	_, err := db.NewCreateTable().
		Model((*Avg)(nil)).
		IfNotExists().
		Exec(context.Background())
	if err != nil {
		return nil, err
	}
	return &AvgRepo{db: db}, nil
}

func (repo *AvgRepo) GetAvgs(ctx context.Context) ([]*Avg, error) {
	var items []*Avg
	err := repo.db.
		NewSelect().
		Model(items).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (repo *AvgRepo) GetAvgByID(ctx context.Context, id string) (*Avg, error) {
	var item Avg
	err := repo.db.
		NewSelect().
		Model(&item).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (repo *AvgRepo) UpsertAvgs(ctx context.Context, avgs []Avg) error {
	return repo.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		for _, avg := range avgs {
			_, err := tx.
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
			if err != nil {
				return err
			}
		}
		return nil
	})
}
