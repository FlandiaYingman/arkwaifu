package avg

import (
	"arkwaifu/internal/app/infra"
	"context"
	"github.com/uptrace/bun"
)

type GroupRepo struct {
	infra.Repo
}

// GroupModel is a group of story. e.g., a 活动 such as "将进酒" or a 主线 such as "怒号光明".
type GroupModel struct {
	bun.BaseModel `bun:"table:groups,alias:g"`

	// ID is the unique id of the AvgGroup.
	// e.g.: "1stact" (骑兵与猎人), "act15side" (将进酒).
	ID string `bun:"id,pk"`

	// Name is the name of the AvgGroup, can be the mainline chapter name, the activity name or the operator record name.
	// e.g.: "骑兵与猎人", "怒号光明", "学者之心", "火山".
	Name string

	// Stories is the stories belong to the group.
	Stories []*StoryModel `bun:"rel:has-many,join:id=group_id"`
}

func NewGroupRepo(db *bun.DB) (*GroupRepo, error) {
	_, err := db.NewCreateTable().
		Model((*GroupModel)(nil)).
		IfNotExists().
		Exec(context.Background())
	if err != nil {
		return nil, err
	}
	return &GroupRepo{
		Repo: infra.NewRepo(db),
	}, nil
}

func (r *GroupRepo) GetGroups(ctx context.Context) ([]GroupModel, error) {
	var items []GroupModel
	err := r.DB.
		NewSelect().
		Model(&items).
		Relation("Stories").
		Relation("Stories.Images").
		Relation("Stories.Backgrounds").
		Scan(ctx)
	return items, err
}

func (r *GroupRepo) GetGroupByID(ctx context.Context, id string) (*GroupModel, error) {
	var item GroupModel
	err := r.DB.
		NewSelect().
		Model(&item).
		Relation("Stories").
		Relation("Stories.Images").
		Relation("Stories.Backgrounds").
		Where("id = ?", id).
		Scan(ctx)
	return &item, err
}

func (r *GroupRepo) InsertGroups(ctx context.Context, groups []GroupModel) error {
	_, err := r.DB.
		NewInsert().
		Model(&groups).
		// On("CONFLICT (id) DO UPDATE").
		// Set("name = EXCLUDED.name").
		Exec(ctx)
	return err
}

func (r *GroupRepo) Truncate(ctx context.Context) error {
	_, err := r.DB.NewTruncateTable().
		Model((*GroupModel)(nil)).
		Exec(ctx)
	return err
}
