package utils

import "fmt"

func FormatTimeElapsed(timeElapsed int64) string {
	if timeElapsed < 0 {
		return "N/A"
	} else if timeElapsed < 1000 {
		return fmt.Sprintf("%dms", timeElapsed)
	} else if timeElapsed < 60000 {
		return fmt.Sprintf("%ds", timeElapsed/1000)
	} else if timeElapsed < 3600000 {
		return fmt.Sprintf("%d:%02dp", timeElapsed/60000, (timeElapsed%60000)/1000)
	}

	return fmt.Sprintf("%02d:%02d:%02d", timeElapsed/3600, (timeElapsed%3600)/60, timeElapsed%60)
}
