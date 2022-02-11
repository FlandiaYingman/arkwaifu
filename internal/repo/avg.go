package repo

import (
	"context"
	"database/sql"
	"github.com/uptrace/bun"
)

// Avg is a fragment of story. e.g.: "8-1 行动前" or "IW-9 行动后" (IW stands for activity "将进酒").
type Avg struct {
	bun.BaseModel `bun:"table:avgs"`

	// StoryID is the unique ID of the Avg.
	// e.g.: "1stact_level_a001_01_beg".
	StoryID string `bun:",pk"`
	// StoryCode is the level code of the Avg, sometimes can be empty (in such case the Avg has no associated level).
	// e.g.: "GT-1", "1-7".
	StoryCode string
	// StoryName is the name of the Avg. The Avg of the same level usually have the same StoryName.
	// e.g.: "埋藏", "我也要大干一场".
	StoryName string
	// StoryTxt is the relative path to the Avg script text, based on "/assets/torappu/dynamicassets/gamedata/" without the extension.
	// e.g.: "activities/act9d0/level_act9d0_06_end" (actually pointing to "/assets/torappu/dynamicassets/gamedata/activities/act9d0/level_act9d0_06_end.txt").
	StoryTxt string
	// AvgTag is the type of the Avg, which can only be "行动前", "行动后" or "幕间".
	AvgTag string

	// GroupID is the ID of the AvgGroup which this Avg belongs to.
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
