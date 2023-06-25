package story

import (
	"github.com/flandiayingman/arkwaifu/internal/app/infra"
	"github.com/flandiayingman/arkwaifu/internal/pkg/ark"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
		Preload("PictureArts").
		Preload("CharacterArts").
		Select("id, server, tag, tag_text, code, name, info, group_id").
		Where("server = ?", server).
		Order("id").
		Find(&stories)
	return stories, result.Error
}
func (r *repo) SelectStory(id string, server ark.Server) (*Story, error) {
	story := new(Story)
	result := r.db.
		Preload("PictureArts").
		Preload("CharacterArts").
		Select("id, server, tag, tag_text, code, name, info, group_id").
		Take(&story, "id = ? AND server = ?", id, server)
	return story, result.Error
}
func (r *repo) SelectStoryGroups(server ark.Server) ([]*Group, error) {
	groups := make([]*Group, 0)
	result := r.db.
		Where("server = ?", server).
		Order("id").
		Preload("Stories", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, server, tag, tag_text, code, name, info, group_id")
		}).
		Preload("Stories.PictureArts").
		Preload("Stories.CharacterArts").
		Find(&groups)
	return groups, result.Error
}
func (r *repo) SelectStoryGroupsByType(server ark.Server, groupType GroupType) ([]*Group, error) {
	groups := make([]*Group, 0)
	result := r.db.
		Where("server = ? AND type = ?", server, groupType).
		Order("id").
		Preload("Stories", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, server, tag, tag_text, code, name, info, group_id")
		}).
		Preload("Stories.PictureArts").
		Preload("Stories.CharacterArts").
		Find(&groups)
	return groups, result.Error
}
func (r *repo) SelectStoryGroup(id string, server ark.Server) (*Group, error) {
	group := new(Group)
	result := r.db.
		Preload("Stories", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, server, tag, tag_text, code, name, info, group_id")
		}).
		Preload("Stories.PictureArts").
		Preload("Stories.CharacterArts").
		Take(&group, "id = ? AND server = ?", id, server)
	return group, result.Error
}

func (r *repo) SelectPictureArts(server ark.Server) ([]*PictureArt, error) {
	arts := make([]*PictureArt, 0)
	return arts, r.db.
		Where("server = ?", server).
		Order("id").
		Find(&arts).Error
}
func (r *repo) SelectPictureArt(server ark.Server, id string) (*PictureArt, error) {
	art := &PictureArt{}
	return art, r.db.
		Where("(server, id) = (?, ?)", server, id).
		Take(&art).Error
}
func (r *repo) SelectCharacterArts(server ark.Server) ([]*CharacterArt, error) {
	arts := make([]*CharacterArt, 0)
	return arts, r.db.
		Where("server = ?", server).
		Order("id").
		Find(&arts).Error
}
func (r *repo) SelectCharacterArt(server ark.Server, id string) (*CharacterArt, error) {
	art := &CharacterArt{}
	return art, r.db.
		Where("(server, id) = (?, ?)", server, id).
		Take(&art).Error
}

func (r *repo) UpsertStoryGroups(groups []Group) error {
	return r.db.
		Clauses(clause.OnConflict{UpdateAll: true}).
		Create(&groups).Error
}