package avg

import (
	"context"
	"github.com/flandiayingman/arkwaifu/internal/app/infra"
	"github.com/uptrace/bun"
)

type StoryRepo struct {
	infra.Repo
}

// storyModel is a part of story of an AVG. e.g., "8-1 行动前" or "IW-9 行动后" (IW stands for activity "将进酒").
type storyModel struct {
	bun.BaseModel `bun:"table:stories,alias:s"`

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

func NewStoryRepo(db *bun.DB) (*StoryRepo, error) {
	_, err := db.NewCreateTable().
		Model((*storyModel)(nil)).
		IfNotExists().
		Exec(context.Background())
	if err != nil {
		return nil, err
	}

	_, err = db.NewCreateTable().
		Model((*assetModel)(nil)).
		IfNotExists().
		Exec(context.Background())
	if err != nil {
		return nil, err
	}
	return &StoryRepo{
		Repo: infra.NewRepo(db),
	}, nil
}

func (r *StoryRepo) GetStories(ctx context.Context) ([]storyModel, error) {
	var items []storyModel
	err := r.DB().
		NewSelect().
		Model(&items).
		Relation("Assets", sortAsset).
		Relation("Group", sortAvg).
		Apply(sortAvg).
		Scan(ctx)
	return items, err
}

func (r *StoryRepo) GetStoryByID(ctx context.Context, id string) (*storyModel, error) {
	var item storyModel
	err := r.DB().
		NewSelect().
		Model(&item).
		Relation("Assets", sortAsset).
		Relation("Group", sortAvg).
		Where("s.id = ?", id).
		Scan(ctx)
	return &item, err
}

func (r *StoryRepo) InsertStories(ctx context.Context, stories []storyModel) error {
	_, err := r.DB().
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
	_, err = r.DB().
		NewInsert().
		Model(&storyToAssets).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

type assetModel struct {
	bun.BaseModel `bun:"table:assets"`
	PK            int64  `bun:"pk,pk,autoincrement"`
	StoryID       string `bun:"story_id"`
	Name          string `bun:"name"`
	Kind          string `bun:"kind"`
}

func (r *StoryRepo) Truncate(ctx context.Context) (err error) {
	err = r.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = r.EndTx(err) }()
	_, err = r.DB().NewTruncateTable().
		Model((*storyModel)(nil)).
		Exec(ctx)
	_, err = r.DB().NewTruncateTable().
		Model((*assetModel)(nil)).
		Exec(ctx)
	return err
}
