package dtos

// GoalItem is the API representation of a category goal, enriched with category name.
type GoalItem struct {
	CategoryId    int    `json:"categoryId"`
	CategoryName  string `json:"categoryName"`
	DailyTargetMs int64  `json:"dailyTargetMs"`
	Enabled       bool   `json:"enabled"`
}

// GoalSet is the input payload for setting a goal.
type GoalSet struct {
	CategoryId    int   `json:"categoryId"`
	DailyTargetMs int64 `json:"dailyTargetMs"`
}
