package data

type CategoryRule struct {
	RuleId     int
	CategoryId CategoryId
	AppName    string
	SubAppName string
}

func (r *CategoryRule) IsMatch(event *AppSwitchEvent) bool {
	if r.SubAppName == "" {
		return r.AppName == event.AppName
	}
	return r.AppName == event.AppName && r.SubAppName == event.SubAppName
}
