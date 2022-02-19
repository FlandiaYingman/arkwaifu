package infra

import (
	"context"
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
	"time"
)

func ProvideFiber() *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      "Ark Waifu",
		ReadTimeout:  time.Second * 20,
		WriteTimeout: time.Second * 20,
	})
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))
	return app
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
