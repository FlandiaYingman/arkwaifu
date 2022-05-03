package infra

import (
	"context"
	"database/sql"
	"time"

	"github.com/flandiayingman/arkwaifu/internal/app/config"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
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

	// db.AddQueryHook(bundebug.NewQueryHook(
	// bundebug.WithVerbose(true),
	// ))
	return db, nil
}
