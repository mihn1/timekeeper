package dtos

import (
	"github.com/mihn1/timekeeper/constants"
	"github.com/mihn1/timekeeper/internal/models"
)

type EventLogItem struct {
	ID           int    `json:"id"`
	AppName      string `json:"appName"`
	StartTime    string `json:"startTime"`   // "15:04:05"
	EndTime      string `json:"endTime"`     // "15:04:05" or "—"
	DurationSecs int64  `json:"durationSecs"`
	CategoryID   int    `json:"categoryId"`
	URLOrTitle   string `json:"urlOrTitle"`  // url if present, else title, else ""
}

func EventLogItemFromModel(e *models.AppSwitchEvent) *EventLogItem {
	url, title := "", ""
	if e.AdditionalData != nil {
		url = e.AdditionalData[constants.KEY_BROWSER_URL]
		title = e.AdditionalData[constants.KEY_BROWSER_TITLE]
	}

	endTime := "—"
	var durationSecs int64
	if !e.EndTime.IsZero() {
		endTime = e.EndTime.Format("15:04:05")
		durationSecs = int64(e.EndTime.Sub(e.StartTime).Seconds())
	}

	urlOrTitle := url
	if urlOrTitle == "" {
		urlOrTitle = title
	}

	return &EventLogItem{
		ID:           int(e.Id),
		AppName:      e.AppName,
		StartTime:    e.StartTime.Format("15:04:05"),
		EndTime:      endTime,
		DurationSecs: durationSecs,
		CategoryID:   int(e.CategoryId),
		URLOrTitle:   urlOrTitle,
	}
}

func EventLogFromModels(events []*models.AppSwitchEvent) []*EventLogItem {
	result := make([]*EventLogItem, len(events))
	for i, e := range events {
		result[i] = EventLogItemFromModel(e)
	}
	return result
}
