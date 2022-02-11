package repo

import (
	"arkwaifu/internal/arkres"
	"context"
	"database/sql"
	"github.com/uptrace/bun"
)

// AvgGroup is a group of Avg. For example, a "活动" such as "将进酒" or a "主线" such as "怒号光明".
type AvgGroup struct {
	bun.BaseModel `bun:"table:avg_groups"`

	ID   string `bun:",pk"`
	Name string
	Avgs []*Avg `bun:"rel:has-many,join:id=group_id"`
}

type AvgGroupRepo struct {
	db *bun.DB
}

func NewAvgGroupRepo(db *bun.DB) (*AvgGroupRepo, error) {
	_, err := db.NewCreateTable().
		Model((*AvgGroup)(nil)).
		IfNotExists().
		Exec(context.Background())
	if err != nil {
		return nil, err
	}
	return &AvgGroupRepo{db: db}, nil
}

func (repo *AvgGroupRepo) GetAvgGroups(ctx context.Context) ([]*AvgGroup, error) {
	var items []*AvgGroup
	err := repo.db.
		NewSelect().
		Model(items).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (repo *AvgGroupRepo) GetAvgGroupByID(ctx context.Context, id string) (*AvgGroup, error) {
	var item AvgGroup
	err := repo.db.
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

func (repo *AvgGroupRepo) UpsertAvgGroups(ctx context.Context, avgGroups []AvgGroup) error {
	return repo.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		for _, group := range avgGroups {
			_, err := tx.
				NewInsert().
				Model(&group).
				On("CONFLICT (id) DO UPDATE").
				Set("name = EXCLUDED.name").
				Exec(ctx)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func AvgGroupFromRaw(raw []arkres.StoryReviewData) []AvgGroup {
	groups := make([]AvgGroup, len(raw))
	for i, d := range raw {
		groups[i] = AvgGroup{
			ID:   d.ID,
			Name: d.Name,
			Avgs: AvgFromRaw(d.InfoUnlockDatas),
		}
	}
	return groups
}

func AvgFromRaw(raw []arkres.StoryData) []*Avg {
	avgs := make([]*Avg, len(raw))
	for i, data := range raw {
		avgs[i] = &Avg{
			StoryID:   data.StoryID,
			StoryCode: data.StoryCode,
			StoryName: data.StoryName,
			StoryTxt:  data.StoryTxt,
			AvgTag:    string(data.AvgTag),
			GroupID:   data.StoryGroup,
		}
	}
	return avgs
}
