package models

type GoalFrequency int

const (
	FrequencyDaily   GoalFrequency = 1
	FrequencyWeekly  GoalFrequency = 2
	FrequencyMonthly GoalFrequency = 3
)

func FrequencyLabel(f GoalFrequency) string {
	switch f {
	case FrequencyWeekly:
		return "Weekly"
	case FrequencyMonthly:
		return "Monthly"
	default:
		return "Daily"
	}
}

type CategoryGoal struct {
	Id          int64
	Name        string
	IsActive    bool
	CategoryIds []CategoryId // in-memory; stored as "1,2,3" in DB
	Frequency   GoalFrequency
	TargetMs    int64
}
