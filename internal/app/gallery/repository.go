package gallery

import (
	"github.com/flandiayingman/arkwaifu/internal/app/infra"
	"github.com/flandiayingman/arkwaifu/internal/pkg/ark"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type repository struct {
	db *infra.Gorm
}

func newRepository(db *infra.Gorm) (*repository, error) {
	r := repository{db: db}
	err := r.init()
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (r *repository) init() (err error) {
	err = r.db.CreateCollateNumeric()
	if err != nil {
		return err
	}
	err = r.db.CreateEnum("game_server",
		ark.CnServer,
		ark.EnServer,
		ark.JpServer,
		ark.KrServer,
		ark.TwServer,
	)
	if err != nil {
		return err
	}
	err = r.db.AutoMigrate(&Gallery{})
	if err != nil {
		return err
	}
	err = r.db.AutoMigrate(&Art{})
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) ListGalleries(server ark.Server) ([]Gallery, error) {
	var galleries []Gallery
	return galleries, r.db.
		Preload("Arts", func(db *gorm.DB) *gorm.DB { return db.Order("gallery_arts.sort_id") }).
		Order("id").
		Find(&galleries, "server = ?", server).
		Error
}

func (r *repository) GetGalleryByID(server ark.Server, id string) (*Gallery, error) {
	var galleries Gallery
	return &galleries, r.db.
		Preload("Arts", func(db *gorm.DB) *gorm.DB { return db.Order("gallery_arts.sort_id") }).
		Take(&galleries, "(server, id) = (?, ?)", server, id).
		Error
}

func (r *repository) ListGalleryArts(server ark.Server) ([]Art, error) {
	var arts []Art
	return arts, r.db.
		Order("gallery_id, sort_id, id").
		Find(&arts, "server = ?", server).
		Error
}

func (r *repository) GetGalleryArtByID(server ark.Server, id string) (*Art, error) {
	var art Art
	return &art, r.db.
		Take(&art, "(server, id) = (?, ?)", server, id).
		Error
}

func (r *repository) Put(g []Gallery) error {
	return r.db.
		Clauses(clause.OnConflict{UpdateAll: true}).
		Create(&g).
		Error
}
