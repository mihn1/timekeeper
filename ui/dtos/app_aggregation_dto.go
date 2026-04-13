package dtos

type AppUsageItem struct {
	AppName      string `json:"appName"`
	TimeElapsed  int64  `json:"timeElapsed"`
	CategoryId   int    `json:"categoryId"`
	CategoryName string `json:"categoryName"`
}
