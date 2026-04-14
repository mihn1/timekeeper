package models

type GoalType string

const (
	GoalTypeDaily   GoalType = "daily"
	GoalTypeWeekly  GoalType = "weekly"
	GoalTypeMonthly GoalType = "monthly"
)

type CategoryGoal struct {
	CategoryId CategoryId
	GoalType   GoalType
	TargetMs   int64
	Enabled    bool
}
