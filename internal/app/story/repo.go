package story

import (
	"github.com/flandiayingman/arkwaifu/internal/app/infra"
	"github.com/flandiayingman/arkwaifu/internal/pkg/ark"
	"gorm.io/gorm"
)

type repo struct {
	db *gorm.DB
}

func (r *repo) init() (err error) {
	err = r.createEnums()
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

func (r *repo) createEnums() error {
	var result *gorm.DB
	result = r.db.Exec(`DO
$$BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'game_server') THEN
        CREATE TYPE game_server AS ENUM ('CN','EN','JP','KR','TW');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'story_tag') THEN
        CREATE TYPE story_tag AS ENUM ('before','after','interlude');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'story_group_type') THEN
        CREATE TYPE story_group_type AS ENUM ('main-story','major-event','minor-event','other');
    END IF;
END$$;`)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func newRepo(db *gorm.DB, _ *infra.NumericCollate) (*repo, error) {
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
func (r *repo) SelectPictureArt(server ark.Server, id string) (*PictureArt, error) {
	art := &PictureArt{}
	return art, r.db.
		Take(&art, "(server, id) = (?, ?)", server, id).Error
}
func (r *repo) SelectCharacterArts(server ark.Server) ([]*CharacterArt, error) {
	arts := make([]*CharacterArt, 0)
	return arts, r.db.
		Order("sort_id").
		Find(&arts, "server = ?", server).Error
}
func (r *repo) SelectCharacterArt(server ark.Server, id string) (*CharacterArt, error) {
	art := &CharacterArt{}
	return art, r.db.
		Take(&art, "(server, id) = (?, ?)", server, id).Error
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
