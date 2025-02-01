package models

type CategoryRule struct {
	RuleId     int
	CategoryId CategoryId
	AppName    string
	Field      string // Path to the field to match with expression
	Expression string
	IsRegex    bool
}

// TODO: Rework this to support additional data for each event
func (r *CategoryRule) IsMatch(event *AppSwitchEvent) bool {
	if r.Field == "" {
		return r.AppName == event.AppName
	}
	return r.AppName == event.AppName
	// && r.Expression == event.SubAppName
}
