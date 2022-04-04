package avg

import (
	"context"
	"github.com/uptrace/bun"
)

// storyModel is a part of story of an AVG. e.g., "8-1 行动前" or "IW-9 行动后" (IW stands for activity "将进酒").
type storyModel struct {
	bun.BaseModel `bun:"table:avg_stories"`

	// ID is the unique ID of the story.
	// e.g., "1stact_level_a001_01_beg".
	ID string `bun:",pk"`

	// Code is the code of the level the story belongs to, could be empty, in this case the story belongs to no level.
	// e.g., "GT-1", "1-7", "".
	Code string

	// Name is the name of the level or the operator record the story belongs to.
	// e.g.: "不要恐慌", "埋藏", "我也要大干一场".
	Name string

	// Tag is the type of the story, which could be only "行动前", "行动后" or "幕间".
	Tag string

	// Images are the images the story uses.
	Assets []*assetModel `bun:"rel:has-many,join:id=story_id"`

	// GroupID is the ID of the group the story belongs to.
	GroupID string
	// Group is the group the story belongs to.
	Group *groupModel `bun:"rel:belongs-to,join:group_id=id"`

	SortID int64 `bun:",autoincrement"`
}
type assetModel struct {
	bun.BaseModel `bun:"table:avg_assets"`
	PK            int64  `bun:"pk,pk,autoincrement"`
	StoryID       string `bun:"story_id"`
	Name          string `bun:"name"`
	Kind          string `bun:"kind"`
}

func (r *Repo) GetStories(ctx context.Context) ([]storyModel, error) {
	var items []storyModel
	err := r.
		NewSelect().
		Model(&items).
		Relation("Assets", sortAsset).
		Relation("Group", sortAvg).
		Apply(sortAvg).
		Scan(ctx)
	return items, err
}
func (r *Repo) GetStoryByID(ctx context.Context, id string) (*storyModel, error) {
	var item storyModel
	err := r.
		NewSelect().
		Model(&item).
		Relation("Assets", sortAsset).
		Relation("Group", sortAvg).
		Where("avg_stories.id = ?", id).
		Scan(ctx)
	return &item, err
}
func (r *Repo) InsertStories(ctx context.Context, stories []storyModel) error {
	_, err := r.
		NewInsert().
		Model(&stories).
		Exec(ctx)
	if err != nil {
		return err
	}

	storyToAssets := make([]assetModel, 0, len(stories))
	for _, story := range stories {
		for _, asset := range story.Assets {
			storyToAssets = append(storyToAssets, *asset)
		}
	}
	_, err = r.
		NewInsert().
		Model(&storyToAssets).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}
func (r *Repo) TruncateStory(ctx context.Context) (err error) {
	return r.DB.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		_, err = r.NewTruncateTable().
			Model((*storyModel)(nil)).
			Exec(ctx)
		if err != nil {
			return err
		}
		_, err = r.NewTruncateTable().
			Model((*assetModel)(nil)).
			Exec(ctx)
		if err != nil {
			return err
		}
		return err
	})
}
