package dtos

// DailyCategorySummary is one (date, category) data point for the multi-day trend chart.
type DailyCategorySummary struct {
	Date         string `json:"date"`
	CategoryId   int    `json:"categoryId"`
	CategoryName string `json:"categoryName"`
	TimeElapsed  int64  `json:"timeElapsed"`
}

// DayActivity is one data point for the year-view heatmap calendar.
type DayActivity struct {
	Date          string `json:"date"`
	TotalMs       int64  `json:"totalMs"`
	TopCategoryId int    `json:"topCategoryId"`
}
