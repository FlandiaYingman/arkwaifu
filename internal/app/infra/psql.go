package infra

import (
	"context"
	"database/sql"
	"github.com/flandiayingman/arkwaifu/internal/app/config"
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
	err := initCustomFunctions(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func initCustomFunctions(db *bun.DB) error {
	// Natural Sort
	// FROM: https://stackoverflow.com/a/48809832/10431637
	// FROM: http://www.rhodiumtoad.org.uk/junk/naturalsort.sql
	_, err := db.DB.Exec(`
		create or replace function natural_sort(text)
		  returns bytea
		  language sql
		  immutable strict
		as $f$
		  	select string_agg(convert_to(coalesce(r[2],length(length(r[1])::text) || length(r[1])::text || r[1]), 'SQL_ASCII'),'\x00')
				from regexp_matches($1, '0*([0-9]+)|([^0-9]+)', 'g') r;
		$f$;
	`)
	if err != nil {
		return err
	}
	return nil
}
