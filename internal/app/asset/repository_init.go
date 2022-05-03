package asset

import (
	"context"
)

func (r *repo) init() error {
	err := r.initTables()
	if err != nil {
		return err
	}
	err = r.initFunctions()
	if err != nil {
		return err
	}
	err = r.initTriggers()
	if err != nil {
		return err
	}
	return nil
}

func (r *repo) initTables() error {
	var err error

	_, err = r.NewCreateTable().
		Model((*modelKindName)(nil)).
		IfNotExists().
		Exec(context.Background())
	if err != nil {
		return err
	}

	_, err = r.NewCreateTable().
		Model((*modelVariantName)(nil)).
		IfNotExists().
		Exec(context.Background())
	if err != nil {
		return err
	}

	_, err = r.NewCreateTable().
		Model((*modelAsset)(nil)).
		IfNotExists().
		ForeignKey("(kind) REFERENCES asset_kind_names (kind_name) ON DELETE CASCADE").
		Exec(context.Background())
	if err != nil {
		return err
	}

	_, err = r.NewCreateTable().
		Model((*modelVariant)(nil)).
		IfNotExists().
		ForeignKey("(variant) REFERENCES asset_variant_names (variant_name) ON DELETE CASCADE").
		ForeignKey("(asset_kind, asset_name) REFERENCES asset_assets (kind, name) ").
		Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) initFunctions() error {
	// Natural Sort
	// FROM: https://stackoverflow.com/a/48809832/10431637
	// FROM: http://www.rhodiumtoad.org.uk/junk/naturalsort.sql
	_, err := r.ExecContext(context.Background(), `
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

func (r *repo) initTriggers() error {
	_, err := r.QueryContext(context.Background(), `
CREATE OR REPLACE FUNCTION update_aa_sort_id()
    RETURNS trigger
    LANGUAGE plpgsql AS
$$
BEGIN
    NEW.kind_sort_id = (SELECT akn.sort_id FROM asset_kind_names AS akn WHERE NEW.kind = akn.kind_name);
    NEW.name_sort_id = (natural_sort(NEW.name));
    return NEW;
END;
$$;

CREATE OR REPLACE FUNCTION update_av_sort_id()
    RETURNS trigger
    LANGUAGE plpgsql AS
$$
BEGIN
    NEW.kind_sort_id = (SELECT aa.kind_sort_id FROM asset_assets AS aa WHERE (NEW.asset_kind, NEW.asset_name) = (aa.kind, aa.name));
    NEW.name_sort_id = (SELECT aa.name_sort_id FROM asset_assets AS aa WHERE (NEW.asset_kind, NEW.asset_name) = (aa.kind, aa.name));
    NEW.variant_sort_id = (SELECT avn.sort_id FROM asset_variant_names AS avn WHERE NEW.variant = avn.variant_name);
    return NEW;
END;
$$;

CREATE OR REPLACE TRIGGER update_aa_sort_id_trigger
    BEFORE INSERT OR UPDATE
    ON asset_assets
    FOR EACH ROW
EXECUTE FUNCTION update_aa_sort_id();

CREATE OR REPLACE TRIGGER update_av_sort_id_trigger
    BEFORE INSERT OR UPDATE
    ON asset_variants
    FOR EACH ROW
EXECUTE FUNCTION update_av_sort_id();
	`)
	return err
}
