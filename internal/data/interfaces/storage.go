package interfaces

type Storage interface {
	Categories() CategoryStore
	Rules() RuleStore
	AppAggregations() AppAggregationStore
	CategoryAggregations() CategoryAggregationStore
	Events() EventStore
	Close() error
}
