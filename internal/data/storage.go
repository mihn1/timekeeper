package data

type Storage interface {
	Categories() CategoryStore
	Rules() RuleStore
	AppAggregations() AppAggregationStore
	CategoryAggregations() CategoryAggregationStore
}
