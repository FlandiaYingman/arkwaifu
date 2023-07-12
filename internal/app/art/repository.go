package art

import (
	"errors"
	"fmt"
	"github.com/flandiayingman/arkwaifu/internal/app/infra"
	"github.com/flandiayingman/arkwaifu/internal/pkg/util/fileutil"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"os"
	"path/filepath"
)

type repository struct {
	db         *gorm.DB
	ContentDir string
}

func newRepo(conf *infra.Config, db *gorm.DB, _ *infra.NumericCollate) (*repository, error) {
	r := repository{
		db:         db,
		ContentDir: filepath.Join(conf.Root, "arts-content"),
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
func (r *repository) SelectArtsByIDs(ids []string) ([]*Art, error) {
	arts := make([]*Art, 0)
	result := r.db.
		Preload("Variants").
		Joins("JOIN UNNEST(?) WITH ORDINALITY t(id, ord) USING (id)", clause.Expr{SQL: "ARRAY[?]", Vars: []any{ids}, WithoutParentheses: true}).
		Order("t.ord").
		Find(&arts, ids)
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

func (r *repository) StoreContent(id string, variation string, content []byte) (err error) {
	err = r.db.Transaction(func(tx *gorm.DB) (err error) {
		// If the corresponding art exists, select it.
		art := Art{}
		err = r.db.
			Where("id = ?", id).
			Take(&art).
			Error
		if err != nil {
			return err
		}

		// If the corresponding variant exists, select it.
		variant := Variant{}
		err = r.db.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("(art_id, variation) = (?, ?)", id, variation).
			Take(&variant).
			Error
		if err != nil {
			return err
		}

		contentObj := Content{
			ID:        art.ID,
			Category:  art.Category,
			Variation: variant.Variation,
			Content:   content,
		}

		// Check whether the content object is valid.
		config, err := contentObj.Check()
		if err != nil {
			return err
		}

		// Update the corresponding variant.
		variant.ContentPresent = true
		variant.ContentWidth = &config.Width
		variant.ContentHeight = &config.Height

		// Write content's change on variant to database
		err = r.UpsertVariants(&variant)
		if err != nil {
			return err
		}

		// Write content to filesystem.
		path := filepath.Join(r.ContentDir, contentObj.PathRel())
		err = fileutil.MkFileFromBytes(path, contentObj.Content)
		if err != nil {
			return err
		}

		return nil
	})
	return
}
func (r *repository) TakeContent(id string, variation string) (content []byte, err error) {
	// If the corresponding art exists, select it.
	art := Art{}
	err = r.db.
		Where("id = ?", id).
		Take(&art).
		Error
	if err != nil {
		return nil, err
	}

	// If the corresponding variant exists, select it.
	variant := Variant{}
	err = r.db.
		Where("(art_id, variation) = (?, ?)", id, variation).
		Take(&variant).
		Error
	if err != nil {
		return nil, err
	}

	if variant.ContentPresent {
		contentObj := Content{
			ID:        art.ID,
			Category:  art.Category,
			Variation: variant.Variation,
			Content:   nil,
		}
		// Read content from the filesystem.
		path := filepath.Join(r.ContentDir, contentObj.PathRel())
		return os.ReadFile(path)
	} else {
		// Error because content is not present.
		return nil, fmt.Errorf("content id=%s variation=%s is not present", id, variation)
	}
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
