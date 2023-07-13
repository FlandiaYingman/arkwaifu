package infra

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"strings"
	"time"
)

type Gorm struct {
	*gorm.DB
}

func ProvideGorm(config *Config) (*Gorm, error) {
	dsn := config.PostgresDSN

	l := logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
		SlowThreshold:             1 * time.Second,
		LogLevel:                  logger.Warn,
		IgnoreRecordNotFoundError: false,
		Colorful:                  true,
	})

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: l,
	})
	if err != nil {
		return nil, err
	}

	return &Gorm{db}, err
}

func (db *Gorm) CreateCollateNumeric() error {
	// The collation is copied from Postgres Docs 24.2.2.3.2. ICU Collations
	// (https://www.postgresql.org/docs/current/collation.html) and 'IF NOT EXISTS' is inserted.
	return db.Exec(`CREATE COLLATION IF NOT EXISTS numeric (provider = icu, locale = 'en-u-kn-true');`).Error
}

func (db *Gorm) CreateEnum(name string, values ...string) error {
	// Surround enum values with double quotes, since it is an identifier.
	name = fmt.Sprintf("\"%s\"", name)
	for i, value := range values {
		// Surround enum values with single quotes, since they are strings
		values[i] = fmt.Sprintf("'%s'", value)
	}
	// Use fmt.Sprintf because a possible BUG occurs when I tried using parameterized query.
	// However, this function is called only by static arguments, so it is considered safe.
	return db.Exec(fmt.Sprintf(
		`DO $$ BEGIN CREATE TYPE %s AS ENUM (%s); EXCEPTION WHEN duplicate_object THEN null; END $$;`,
		name,
		strings.Join(values, ","),
	)).Error
}
