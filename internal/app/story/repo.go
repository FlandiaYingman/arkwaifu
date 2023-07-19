package story

import (
	"github.com/flandiayingman/arkwaifu/internal/app/infra"
	"github.com/flandiayingman/arkwaifu/internal/pkg/ark"
	"gorm.io/gorm"
)

type repo struct {
	db *infra.Gorm
}

func (r *repo) init() (err error) {
	err = r.db.CreateEnum("game_server", "CN", "EN", "JP", "KR", "TW")
	if err != nil {
		return err
	}
	err = r.db.CreateEnum("story_tag", "before", "after", "interlude")
	if err != nil {
		return err
	}
	err = r.db.CreateEnum("story_group_type", "main-story", "major-event", "minor-event", "other")
	if err != nil {
		return err
	}
	err = r.db.AutoMigrate(&Group{})
	if err != nil {
		return err
	}
	err = r.db.AutoMigrate(&Story{})
	if err != nil {
		return err
	}
	err = r.db.AutoMigrate(&PictureArt{})
	if err != nil {
		return err
	}
	err = r.db.AutoMigrate(&CharacterArt{})
	if err != nil {
		return err
	}
	return nil
}

func newRepo(db *infra.Gorm) (*repo, error) {
	r := repo{db}
	err := r.init()
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (r *repo) SelectStories(server ark.Server) ([]*Story, error) {
	stories := make([]*Story, 0)
	result := r.db.
		Preload("PictureArts", func(db *gorm.DB) *gorm.DB { return db.Order("picture_arts.sort_id") }).
		Preload("CharacterArts", func(db *gorm.DB) *gorm.DB { return db.Order("character_arts.sort_id") }).
		Order("sort_id").
		Find(&stories, "server = ?", server)
	return stories, result.Error
}
func (r *repo) SelectStory(id string, server ark.Server) (*Story, error) {
	story := &Story{}
	result := r.db.
		Preload("PictureArts", func(db *gorm.DB) *gorm.DB { return db.Order("picture_arts.sort_id") }).
		Preload("CharacterArts", func(db *gorm.DB) *gorm.DB { return db.Order("character_arts.sort_id") }).
		Take(&story, "(id, server) = (?, ?)", id, server)
	return story, result.Error
}
func (r *repo) SelectStoryGroups(server ark.Server) ([]*Group, error) {
	groups := make([]*Group, 0)
	result := r.db.
		Preload("Stories", func(db *gorm.DB) *gorm.DB { return db.Order("stories.sort_id") }).
		Preload("Stories.PictureArts", func(db *gorm.DB) *gorm.DB { return db.Order("picture_arts.sort_id") }).
		Preload("Stories.CharacterArts", func(db *gorm.DB) *gorm.DB { return db.Order("character_arts.sort_id") }).
		Order("sort_id").
		Find(&groups, "server = ?", server)
	return groups, result.Error
}
func (r *repo) SelectStoryGroupsByType(server ark.Server, groupType GroupType) ([]*Group, error) {
	groups := make([]*Group, 0)
	result := r.db.
		Preload("Stories", func(db *gorm.DB) *gorm.DB { return db.Order("stories.sort_id") }).
		Preload("Stories.PictureArts", func(db *gorm.DB) *gorm.DB { return db.Order("picture_arts.sort_id") }).
		Preload("Stories.CharacterArts", func(db *gorm.DB) *gorm.DB { return db.Order("character_arts.sort_id") }).
		Order("sort_id").
		Find(&groups, "(server, type) = (?, ?)", server, groupType)
	return groups, result.Error
}
func (r *repo) SelectStoryGroup(id string, server ark.Server) (*Group, error) {
	group := &Group{}
	result := r.db.
		Preload("Stories", func(db *gorm.DB) *gorm.DB { return db.Order("stories.sort_id") }).
		Preload("Stories.PictureArts", func(db *gorm.DB) *gorm.DB { return db.Order("picture_arts.sort_id") }).
		Preload("Stories.CharacterArts", func(db *gorm.DB) *gorm.DB { return db.Order("character_arts.sort_id") }).
		Take(&group, "(id, server) = (?, ?)", id, server)
	return group, result.Error
}

func (r *repo) SelectPictureArts(server ark.Server) ([]*PictureArt, error) {
	arts := make([]*PictureArt, 0)
	return arts, r.db.
		Order("sort_id").
		Find(&arts, "server = ?", server).Error
}
func (r *repo) SelectAggregatedPictureArtByID(server ark.Server, id string) (*AggregatedPictureArt, error) {
	art := &AggregatedPictureArt{}
	err := r.db.
		Model(&PictureArt{}).
		Take(&art, "(server, id) = (?, ?)", server, id).
		Error
	if err != nil {
		return nil, err
	}
	return art, nil
}
func (r *repo) SelectCharacterArts(server ark.Server) ([]*CharacterArt, error) {
	arts := make([]*CharacterArt, 0)
	return arts, r.db.
		Order("sort_id").
		Find(&arts, "server = ?", server).Error
}
func (r *repo) SelectAggregatedCharacterArtByID(server ark.Server, id string) (*AggregatedCharacterArt, error) {
	var err error
	art := &AggregatedCharacterArt{}

	err = r.db.
		Model(&CharacterArt{}).
		Take(&art, "(server, id) = (?, ?)", server, id).
		Error
	if err != nil {
		return nil, err
	}

	err = r.db.
		Table("(?) as dt(aggregated_names)",
			r.db.
				Model(&CharacterArt{}).
				Select("unnest(names)").
				Where("(server, id) = (?, ?)", server, id),
		).
		Select("array_agg(DISTINCT aggregated_names) as names").
		Scan(&art).
		Error
	if err != nil {
		return nil, err
	}

	return art, nil
}

func (r *repo) UpsertStoryGroups(groups []Group) error {
	// UpsertStoryGroups does not actually upsert values.
	// Because of the existence of SortID (we want to re-generate it), values are deleted and then inserted.
	return r.db.Transaction(func(tx *gorm.DB) error {
		var err error
		err = tx.Delete(&groups).Error
		if err != nil {
			return err
		}
		err = tx.Create(&groups).Error
		if err != nil {
			return err
		}
		return nil
	})
}
