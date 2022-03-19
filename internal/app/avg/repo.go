package avg

import "github.com/uptrace/bun"

func sortAvg(query *bun.SelectQuery) *bun.SelectQuery {
	return query.Order("sort_id ASC")
}
func sortAsset(query *bun.SelectQuery) *bun.SelectQuery {
	return query.Order("pk ASC")
}
