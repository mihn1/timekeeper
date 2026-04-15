package dtos

// GoalItem is the API representation of a goal, enriched with category names.
type GoalItem struct {
	Id            int64    `json:"id"`
	Name          string   `json:"name"`
	IsActive      bool     `json:"isActive"`
	CategoryIds   []int    `json:"categoryIds"`
	CategoryNames []string `json:"categoryNames"`
	Frequency     int      `json:"frequency"`
	TargetMs      int64    `json:"targetMs"`
}
