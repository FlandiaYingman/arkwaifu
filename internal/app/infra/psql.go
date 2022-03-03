package infra

import (
	"context"
	"database/sql"
	"github.com/flandiayingman/arkwaifu/internal/app/config"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"time"
)

func ProvidePostgres(config *config.Config) (*bun.DB, error) {
	opts := []pgdriver.Option{
		pgdriver.WithDSN(config.PostgresDSN),
		pgdriver.WithApplicationName("arkwaifu"),
	}

	pgdb := sql.OpenDB(pgdriver.NewConnector(opts...))
	db := bun.NewDB(pgdb, pgdialect.New())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}
	return db, nil
}

type Repo struct {
	// conn is for injection. e.g., Atomic.
	conn *bun.DB
	// DB is for normal database operations.
	DB bun.IDB
}

func NewRepo(db *bun.DB) Repo {
	return Repo{
		conn: db,
		DB:   db,
	}
}

func (r *Repo) BeginTx(ctx context.Context) (err error) {
	tx, err := r.conn.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	r.DB = tx
	return nil
}

func (r *Repo) Commit() error {
	tx, ok := r.DB.(bun.Tx)
	if !ok {
		return errors.New("The transaction hasn't begun.")
	}
	err := tx.Commit()
	if err != nil {
		return err
	}
	r.DB = r.conn
	return nil
}

func (r *Repo) Rollback() error {
	tx, ok := r.DB.(bun.Tx)
	if !ok {
		return errors.New("The transaction hasn't begun.")
	}
	err := tx.Rollback()
	if err != nil {
		return err
	}
	r.DB = r.conn
	return nil
}

func (r *Repo) EndTx(err error) error {
	if err != nil {
		return r.Rollback()
	} else {
		return r.Commit()
	}
}
