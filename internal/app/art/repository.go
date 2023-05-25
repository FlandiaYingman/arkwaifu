package art

import (
	"errors"
	"github.com/flandiayingman/arkwaifu/internal/app/infra"
	"github.com/flandiayingman/arkwaifu/internal/pkg/util/fileutil"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"os"
)

type repository struct {
	db        *gorm.DB
	staticDir string
}

func newRepo(conf *infra.Config, db *gorm.DB, _ *infra.NumericCollate) (*repository, error) {
	r := repository{
		db:        db,
		staticDir: conf.Root,
	}
	err := r.init()
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (r *repository) init() (err error) {
	err = r.db.AutoMigrate(&Art{})
	if err != nil {
		return err
	}
	err = r.db.AutoMigrate(&Variant{})
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) SelectArts() ([]*Art, error) {
	arts := make([]*Art, 0)
	result := r.db.
		Preload("Variants").
		Order("id").
		Find(&arts)
	if result.Error != nil {
		return nil, result.Error
	}
	return arts, nil
}
func (r *repository) SelectArtsByCategory(category string) ([]*Art, error) {
	arts := make([]*Art, 0)
	result := r.db.
		Preload("Variants").
		Where("category = ?", category).
		Order("id").
		Find(&arts)
	if result.Error != nil {
		return nil, result.Error
	}
	return arts, nil
}
func (r *repository) SelectArt(id string) (*Art, error) {
	art := new(Art)
	result := r.db.Preload("Variants").Take(&art, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.Join(ErrNotFound, result.Error)
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return art, nil
}

func (r *repository) SelectVariants(id string) ([]*Variant, error) {
	variants := make([]*Variant, 0)
	result := r.db.
		Where("art_id = ?", id).
		Order("art_id").
		Order("kind").
		Find(&variants)
	if result.Error != nil {
		return nil, result.Error
	}
	return variants, nil
}
func (r *repository) SelectVariant(artID string, typ string) (*Variant, error) {
	variant := new(Variant)
	result := r.db.Take(&variant, "art_id = ?, type = ?", artID, typ)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.Join(ErrNotFound, result.Error)
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return variant, nil
}

func (r *repository) UpsertArts(arts ...*Art) error {
	return r.db.
		Clauses(clause.OnConflict{UpdateAll: true}).
		Create(&arts).Error
}
func (r *repository) UpsertVariants(variants ...*Variant) error {
	err := r.db.
		Clauses(clause.OnConflict{UpdateAll: true}).
		Create(&variants).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) StoreStatics(statics ...*VariantContent) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, static := range statics {
			var err error

			// Select the corresponding variant and check whether it exists
			v := Variant{}
			err = r.db.
				Clauses(clause.Locking{Strength: "UPDATE"}).
				Where("(art_id, variation) = (?, ?)", static.ArtID, static.Variation).
				Take(&v).
				Error
			if err != nil {
				return err
			}

			// Check whether static is valid
			config, err := static.Check()
			if err != nil {
				return err
			}

			// Update the corresponding variant
			v.ContentPresent = true
			v.ContentPath = static.PathRel()
			v.ContentWidth, v.ContentHeight = &config.Width, &config.Height

			// Write to file system
			path := static.Path(r.staticDir)
			err = fileutil.MkFileFromBytes(path, static.Content)
			if err != nil {
				return err
			}

			// Write change to database
			err = r.UpsertVariants(&v)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
func (r *repository) TakeStatics(statics ...*VariantContent) error {
	var err error
	for _, static := range statics {
		// Select the corresponding variant and check whether it exists
		v := Variant{}
		err = r.db.
			Model(&Variant{}).
			Select("").
			Where("(art_id, variation) = (?, ?)", static.ArtID, static.Variation).
			Take(&v).
			Error
		if err != nil {
			return err
		}

		// Read the static file
		static.Content, err = os.ReadFile(static.Path(r.staticDir))
		if err != nil {
			return err
		}

		// Check whether static is valid
		_, err = static.Check()
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *repository) SelectArtsWhereVariantAbsent(variation string) ([]*Art, error) {
	arts := make([]*Art, 0)
	err := r.db.
		Preload("Variants").
		Order("id").
		Not("EXISTS (SELECT 1 FROM variants WHERE (arts.id, ?) = (variants.art_id, variants.variation))", variation).
		Find(&arts).
		Error
	return arts, err
}
