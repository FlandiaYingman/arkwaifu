package entity

import (
	"github.com/uptrace/bun"
)

// Avg is a fragment of story. e.g.: "8-1 行动前" or "IW-9 行动后" (IW stands for activity "将进酒").
type Avg struct {
	bun.BaseModel `bun:"table:avgs"`

	// StoryID is the unique ID of the Avg.
	// e.g.: "1stact_level_a001_01_beg".
	StoryID string `bun:",pk"`
	// StoryCode is the level code of the Avg, sometimes can be empty (in such case the Avg has no associated level).
	// e.g.: "GT-1", "1-7".
	StoryCode string
	// StoryName is the name of the Avg. The Avg of the same level usually have the same StoryName.
	// e.g.: "埋藏", "我也要大干一场".
	StoryName string
	// StoryTxt is the Avg script text.
	// e.g.: "activities/act9d0/level_act9d0_06_end" (actually pointing to "/assets/torappu/dynamicassets/gamedata/activities/act9d0/level_act9d0_06_end.txt").
	StoryTxt string
	// AvgTag is the type of the Avg, which can only be "行动前", "行动后" or "幕间".
	AvgTag string

	// GroupID is the ID of the AvgGroup which this Avg belongs to.
	GroupID string `bun:"group_id"`
}

// AvgGroup is a group of Avg. For example, a "活动" such as "将进酒" or a "主线" such as "怒号光明".
type AvgGroup struct {
	bun.BaseModel `bun:"table:avg_groups"`

	// ID is the unique id of the AvgGroup.
	// e.g.: "1stact" (骑兵与猎人), "act15side" (将进酒).
	ID string `bun:",pk"`
	// Name is the name of the AvgGroup, can be the mainline chapter name, the activity name or the operator record name.
	// e.g.: "骑兵与猎人", "怒号光明", "学者之心", "火山".
	Name string

	// Avgs is the Avg that belong to the AvgGroup. Currently, each Avg belongs to one and only one AvgGroup.
	Avgs []*Avg `bun:"rel:has-many,join:id=group_id"`
}
