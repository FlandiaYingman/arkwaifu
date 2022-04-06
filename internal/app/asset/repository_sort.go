package asset

import "github.com/uptrace/bun"

func SortAsset(query *bun.SelectQuery) *bun.SelectQuery {
	return query.
		OrderExpr("(SELECT asset_kind_names.sort_id FROM asset_kind_names WHERE kind = asset_kind_names.kind_name)").
		OrderExpr("(natural_sort(name))")
}

func SortVariant(query *bun.SelectQuery) *bun.SelectQuery {
	return query.
		OrderExpr("(SELECT asset_kind_names.sort_id FROM asset_kind_names WHERE asset_kind = asset_kind_names.kind_name)").
		OrderExpr("(natural_sort(asset_name))").
		OrderExpr("(SELECT asset_variant_names.sort_id FROM asset_variant_names WHERE variant = asset_variant_names.variant_name)")
}
