package asset

import (
	"context"
	"io"

	"github.com/flandiayingman/arkwaifu/internal/app/config"
	"github.com/flandiayingman/arkwaifu/internal/pkg/util/fileutil"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/uptrace/bun"
)

type repo struct {
	bun.IDB
	DB *bun.DB

	StaticDir string
}

func NewRepo(db *bun.DB, conf *config.Config) (*repo, error) {
	r := repo{
		IDB:       db,
		DB:        db,
		StaticDir: conf.StaticDir,
	}

	err := r.init()
	if err != nil {
		return nil, err
	}

	return &r, nil
}

var (
	ErrExists = errors.New("asset: the asset or variant already exists.")
)

func (r *repo) Truncate(ctx context.Context) error {
	return r.DB.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		var err error
		_, err = tx.
			NewTruncateTable().
			Model((*mAsset)(nil)).
			Cascade().
			Exec(ctx)
		if err != nil {
			return err
		}
		_, err = tx.
			NewTruncateTable().
			Model((*mVariant)(nil)).
			Cascade().
			Exec(ctx)
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *repo) InsertAsset(ctx context.Context, models ...mAsset) error {
	if len(models) == 0 {
		return nil
	}
	_, err := r.
		NewInsert().
		Model(&models).
		Exec(ctx)
	return err
}
func (r *repo) InsertVariant(ctx context.Context, models ...mVariant) error {
	if len(models) == 0 {
		return nil
	}
	_, err := r.
		NewInsert().
		Model(&models).
		Exec(ctx)
	return err
}
func (r *repo) InsertVariantFile(ctx context.Context, m mVariant, f io.Reader) error {
	return r.DB.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		// Insert the variant. The variant is invisible until the transaction is committed.
		res, err := tx.NewInsert().
			Model(&m).
			On("CONFLICT DO NOTHING").
			Exec(ctx)
		if err != nil {
			return err
		}

		// Check if the insertion conflicts. If it conflicts, return ErrExists.
		rows, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if rows <= 0 {
			return ErrExists
		}

		// Write variant file. Note MkFileFromReader closes the reader.
		err = fileutil.MkFileFromReader(m.FilePath(r.StaticDir), f)
		if err != nil {
			return err
		}

		// Return nil to commit if everything success.
		return nil
	})
}

func (r *repo) SelectAssets(ctx context.Context, kind string) ([]mAsset, error) {
	models := new([]mAsset)
	query := r.
		NewSelect().
		Model(models).
		Relation("Variants", SortAssetVariant).
		Apply(SortAsset)
	if kind != "" {
		query.Where("kind = ?", kind)
	}
	err := query.Scan(ctx)
	return *models, err
}
func (r *repo) SelectAsset(ctx context.Context, kind, name string) (*mAsset, error) {
	model := new(mAsset)
	err := r.
		NewSelect().
		Model(model).
		Relation("Variants", SortAssetVariant).
		Where("(kind, name) = (?, ?)", kind, name).
		Scan(ctx)
	return model, err
}

func (r *repo) SelectVariants(ctx context.Context, kind, name string) ([]mVariant, error) {
	models := new([]mVariant)
	err := r.
		NewSelect().
		Model(models).
		Where("(asset_kind, asset_name) = (?, ?)", kind, name).
		Apply(SortVariant).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return *models, err
}
func (r *repo) SelectVariant(ctx context.Context, kind, name, variant string) (*mVariant, error) {
	model := new(mVariant)
	err := r.
		NewSelect().
		Model(model).
		Where("(asset_kind, asset_name, variant) = (?, ?, ?)", kind, name, variant).
		Scan(ctx)
	return model, err
}

func (r *repo) InitNames(ctx context.Context, kindNames []string, variantNames []string) error {
	kms := lo.Map(kindNames, func(s string, _ int) mKindName { return mKindName{KindName: s} })
	vms := lo.Map(variantNames, func(s string, _ int) mVariantName { return mVariantName{VariantName: s} })
	return r.DB.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		r.IDB = tx
		var err error
		_, err = tx.
			NewTruncateTable().
			Model((*mKindName)(nil)).
			Cascade().
			Exec(ctx)
		if err != nil {
			return err
		}
		_, err = tx.
			NewTruncateTable().
			Model((*mVariantName)(nil)).
			Cascade().
			Exec(ctx)
		if err != nil {
			return err
		}
		_, err = r.
			NewInsert().
			Model(&kms).
			Exec(ctx)
		if err != nil {
			return err
		}
		_, err = r.
			NewInsert().
			Model(&vms).
			Exec(ctx)
		if err != nil {
			return err
		}
		return nil
	})
}
func (r *repo) SelectKindNames(ctx context.Context) ([]string, error) {
	models := new([]mKindName)
	err := r.
		NewSelect().
		Model(models).
		Order("sort_id").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return lo.Map(*models, func(m mKindName, _ int) string { return m.KindName }), nil
}
func (r *repo) SelectVariantNames(ctx context.Context) ([]string, error) {
	models := new([]mVariantName)
	err := r.
		NewSelect().
		Model(models).
		Order("sort_id").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return lo.Map(*models, func(m mVariantName, _ int) string { return m.VariantName }), nil
}

func (r repo) Update(ctx context.Context, ams []mAsset, vms []mVariant) error {
	return r.DB.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		r.IDB = tx
		var err error
		err = r.Truncate(ctx)
		if err != nil {
			return errors.Wrapf(err, "failed to truncate asset table")
		}
		err = r.InsertAsset(ctx, ams...)
		if err != nil {
			return errors.Wrapf(err, "failed to insert assets %v", ams)
		}
		err = r.InsertVariant(ctx, vms...)
		if err != nil {
			return errors.Wrapf(err, "failed to insert variants %v", vms)
		}
		return nil
	})
}
