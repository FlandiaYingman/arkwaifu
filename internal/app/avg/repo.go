package avg

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

type Repo struct {
	bun.IDB
	DB *bun.DB
}

type versionModel struct {
	bun.BaseModel `bun:"table:version"`

	// ID can only be true, because this table should only have one row.
	ID         *bool  `bun:",pk,default:true"`
	ResVersion string `bun:""`
}

func (r *Repo) GetVersion(ctx context.Context) (string, error) {
	var resVersion versionModel
	err := r.
		NewSelect().
		Model(&resVersion).
		Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return resVersion.ResVersion, nil
}
func (r *Repo) UpsertVersion(ctx context.Context, resVersion string) error {
	resVersionEntity := versionModel{
		ResVersion: resVersion,
	}
	_, err := r.
		NewInsert().
		Model(&resVersionEntity).
		On("CONFLICT (id) DO UPDATE").
		Set("res_version = EXCLUDED.res_version").
		Exec(ctx)
	return err
}

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

// groupModel is a group of story. e.g., a 活动 such as "将进酒" or a 主线 such as "怒号光明".
type groupModel struct {
	bun.BaseModel `bun:"table:avg_groups"`

	// ID is the unique id of the AvgGroup.
	// e.g.: "1stact" (骑兵与猎人), "act15side" (将进酒).
	ID string `bun:"id,pk"`

	// Name is the name of the AvgGroup, can be the mainline chapter name, the activity name or the operator record name.
	// e.g.: "骑兵与猎人", "怒号光明", "学者之心", "火山".
	Name string

	Type string

	// Stories is the stories belong to the group.
	Stories []*storyModel `bun:"rel:has-many,join:id=group_id"`
	SortID  int64         `bun:",autoincrement"`
}

func (r *Repo) GetGroups(ctx context.Context) ([]groupModel, error) {
	var items []groupModel
	err := r.
		NewSelect().
		Model(&items).
		Relation("Stories", sortAvg).
		Relation("Stories.Assets", sortAsset).
		Apply(sortAvg).
		Scan(ctx)
	return items, err
}
func (r *Repo) GetGroupByID(ctx context.Context, id string) (*groupModel, error) {
	var item groupModel
	err := r.
		NewSelect().
		Model(&item).
		Relation("Stories", sortAvg).
		Relation("Stories.Assets", sortAsset).
		Where("avg_groups.id = ?", id).
		Scan(ctx)
	return &item, err
}
func (r *Repo) InsertGroups(ctx context.Context, groups []groupModel) error {
	_, err := r.
		NewInsert().
		Model(&groups).
		// On("CONFLICT (id) DO UPDATE").
		// Set("name = EXCLUDED.name").
		Exec(ctx)
	return err
}
func (r *Repo) Truncate(ctx context.Context) error {
	_, err := r.NewTruncateTable().
		Model((*groupModel)(nil)).
		Exec(ctx)
	return err
}

func NewRepo(db *bun.DB) (*Repo, error) {
	r := Repo{DB: db, IDB: db}
	err := r.DB.RunInTx(context.Background(), nil, func(ctx context.Context, tx bun.Tx) error {
		var err error
		_, err = db.NewCreateTable().
			Model((*versionModel)(nil)).
			IfNotExists().
			ColumnExpr("CHECK (id = ?)", true).
			Exec(context.Background())
		if err != nil {
			return err
		}
		_, err = db.NewCreateTable().
			Model((*groupModel)(nil)).
			IfNotExists().
			Exec(context.Background())
		if err != nil {
			return err
		}
		_, err = db.NewCreateTable().
			Model((*storyModel)(nil)).
			IfNotExists().
			Exec(context.Background())
		if err != nil {
			return err
		}
		_, err = db.NewCreateTable().
			Model((*assetModel)(nil)).
			IfNotExists().
			Exec(context.Background())
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (r Repo) UpdateAvg(ctx context.Context, version string, gms []groupModel, sms []storyModel) error {
	return r.DB.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		r.IDB = tx
		var err error

		err = r.Truncate(ctx)
		if err != nil {
			return err
		}
		err = r.Truncate(ctx)
		if err != nil {
			return err
		}
		err = r.UpsertVersion(ctx, version)
		if err != nil {
			return err
		}

		err = r.InsertGroups(ctx, gms)
		if err != nil {
			return err
		}

		err = r.InsertStories(ctx, sms)
		if err != nil {
			return err
		}

		return nil
	})
}
