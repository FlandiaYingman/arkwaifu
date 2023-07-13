package updateloop

import (
	"github.com/flandiayingman/arkwaifu/internal/app/infra"
	"github.com/flandiayingman/arkwaifu/internal/pkg/ark"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type artVersion struct {
	Lock    *int        `gorm:"primaryKey;check:lock=0"`
	Version ark.Version `gorm:""`
}

type storyVersion struct {
	Server  ark.Server  `gorm:"primaryKey;type:game_server"`
	Version ark.Version `gorm:""`
}

type repo struct {
	*infra.Gorm
}

var zeroPtr = func() *int {
	zero := 0
	return &zero
}()

func newRepo(db *infra.Gorm) (*repo, error) {
	var err error
	repo := &repo{db}
	err = repo.initArtVersionTable()
	if err != nil {
		return nil, err
	}
	err = repo.initStoryVersionTable()
	if err != nil {
		return nil, err
	}
	return repo, nil
}

func (r *repo) initArtVersionTable() error {
	err := r.AutoMigrate(&artVersion{})
	if err != nil {
		return err
	}
	var result *gorm.DB
	result = r.
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&artVersion{
			Lock:    zeroPtr,
			Version: "",
		})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *repo) initStoryVersionTable() error {
	err := r.AutoMigrate(&storyVersion{})
	if err != nil {
		return err
	}
	var result *gorm.DB
	result = r.
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&storyVersion{Server: ark.CnServer})
	if result.Error != nil {
		return result.Error
	}
	result = r.
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&storyVersion{Server: ark.EnServer})
	if result.Error != nil {
		return result.Error
	}
	result = r.
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&storyVersion{Server: ark.JpServer})
	if result.Error != nil {
		return result.Error
	}
	result = r.
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&storyVersion{Server: ark.KrServer})
	if result.Error != nil {
		return result.Error
	}
	result = r.
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&storyVersion{Server: ark.TwServer})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *repo) upsertArtVersion(v *artVersion) error {
	return r.
		Clauses(clause.OnConflict{UpdateAll: true}).
		Create(&v).
		Error
}
func (r *repo) upsertStoryVersion(v *storyVersion) error {
	return r.
		Clauses(clause.OnConflict{UpdateAll: true}).
		Create(&v).
		Error
}
func (r *repo) selectArtVersion() (ark.Version, error) {
	var v artVersion
	err := r.Take(&v).Error
	return v.Version, err
}
func (r *repo) selectStoryVersion(server ark.Server) (ark.Version, error) {
	var v storyVersion
	err := r.
		Where("server = ?", server).
		Take(&v).Error
	return v.Version, err
}
