package asset

import "github.com/uptrace/bun"

func SortAsset(query *bun.SelectQuery) *bun.SelectQuery {
	return query.Order("kind_sort_id", "name_sort_id")
}

func SortVariant(query *bun.SelectQuery) *bun.SelectQuery {
	return query.Order("kind_sort_id", "name_sort_id", "variant_sort_id")
}

func SortAssetVariant(query *bun.SelectQuery) *bun.SelectQuery {
	return query.Order("variant_sort_id")
}
