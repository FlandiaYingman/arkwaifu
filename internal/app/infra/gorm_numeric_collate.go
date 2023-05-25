package infra

import "gorm.io/gorm"

type NumericCollate struct {
}

func ProvideGormNumericCollate(db *gorm.DB) (*NumericCollate, error) {
	// The collation is copied from Postgres Docs 24.2.2.3.2. ICU Collations
	// (https://www.postgresql.org/docs/current/collation.html) and 'IF NOT EXISTS' is inserted.
	result := db.Exec(`CREATE COLLATION IF NOT EXISTS numeric (provider = icu, locale = 'en-u-kn-true');`)
	return &NumericCollate{}, result.Error
}
