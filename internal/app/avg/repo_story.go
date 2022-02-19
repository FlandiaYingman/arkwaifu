package avg

import (
	"arkwaifu/internal/app/infra"
	"context"
	"github.com/uptrace/bun"
)

type StoryRepo struct {
	infra.Repo
}

// StoryModel is a part of story of an AVG. e.g., "8-1 行动前" or "IW-9 行动后" (IW stands for activity "将进酒").
type StoryModel struct {
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
	Images []*ImageModel `bun:"rel:has-many,join:id=story_id"`

	// Backgrounds are the backgrounds the story uses.
	Backgrounds []*BackgroundModel `bun:"rel:has-many,join:id=story_id"`

	// GroupID is the ID of the group the story belongs to.
	GroupID string
	// Group is the group the story belongs to.
	Group *GroupModel `bun:"rel:belongs-to,join:group_id=id"`
}

func NewStoryRepo(db *bun.DB) (*StoryRepo, error) {
	_, err := db.NewCreateTable().
		Model((*StoryModel)(nil)).
		IfNotExists().
		Exec(context.Background())
	if err != nil {
		return nil, err
	}

	_, err = db.NewCreateTable().
		Model((*ImageModel)(nil)).
		IfNotExists().
		Exec(context.Background())
	if err != nil {
		return nil, err
	}

	_, err = db.NewCreateTable().
		Model((*BackgroundModel)(nil)).
		IfNotExists().
		Exec(context.Background())
	if err != nil {
		return nil, err
	}
	return &StoryRepo{
		Repo: infra.NewRepo(db),
	}, nil
}

func (r *StoryRepo) GetStories(ctx context.Context) ([]StoryModel, error) {
	var items []StoryModel
	err := r.DB.
		NewSelect().
		Model(&items).
		Relation("Images").
		Relation("Backgrounds").
		Relation("Group").
		Scan(ctx)
	return items, err
}

func (r *StoryRepo) GetStoryByID(ctx context.Context, id string) (*StoryModel, error) {
	var item StoryModel
	err := r.DB.
		NewSelect().
		Model(&item).
		Relation("Images").
		Relation("Backgrounds").
		Relation("Group").
		Where("id = ?", id).
		Scan(ctx)
	return &item, err
}

func (r *StoryRepo) InsertStories(ctx context.Context, stories []StoryModel) error {
	_, err := r.DB.
		NewInsert().
		Model(&stories).
		// On("CONFLICT (id) DO UPDATE").
		// Set("id = EXCLUDED.id").
		// Set("code = EXCLUDED.code").
		// Set("name = EXCLUDED.name").
		// Set("tag = EXCLUDED.tag").
		// Set("group_id = EXCLUDED.group_id").
		Exec(ctx)
	if err != nil {
		return err
	}

	storyToImages := make([]ImageModel, 0, len(stories))
	storyToBackgrounds := make([]BackgroundModel, 0, len(stories))
	for _, story := range stories {
		for _, image := range story.Images {
			storyToImages = append(storyToImages, *image)
		}
		for _, background := range story.Backgrounds {
			storyToBackgrounds = append(storyToBackgrounds, *background)
		}
	}
	_, err = r.DB.
		NewInsert().
		Model(&storyToImages).
		Exec(ctx)
	if err != nil {
		return err
	}
	_, err = r.DB.
		NewInsert().
		Model(&storyToBackgrounds).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

type ImageModel struct {
	bun.BaseModel `bun:"table:images"`
	ID            int64  `bun:"id,pk,autoincrement"`
	StoryID       string `bun:""`
	Image         string `bun:""`
}

type BackgroundModel struct {
	bun.BaseModel `bun:"table:backgrounds"`
	ID            int64  `bun:"id,pk,autoincrement"`
	StoryID       string `bun:""`
	Background    string `bun:""`
}

func (r *StoryRepo) Truncate(ctx context.Context) (err error) {
	err = r.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = r.EndTx(err) }()
	_, err = r.DB.NewTruncateTable().
		Model((*StoryModel)(nil)).
		Exec(ctx)
	_, err = r.DB.NewTruncateTable().
		Model((*ImageModel)(nil)).
		Exec(ctx)
	_, err = r.DB.NewTruncateTable().
		Model((*BackgroundModel)(nil)).
		Exec(ctx)
	return err
}
