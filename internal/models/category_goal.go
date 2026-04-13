package models

type CategoryGoal struct {
	CategoryId    CategoryId
	DailyTargetMs int64
	Enabled       bool
}
