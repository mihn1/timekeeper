package data

type Storage struct {
	CategoryStore            CategoryStore
	RuleStore                RuleStore
	AppAggregationStore      AppAggregationStore
	CategoryAggregationStore CategoryAggregationStore
}

func NewStorage(
	categoryStore CategoryStore,
	ruleStore RuleStore,
	appAggregationStore AppAggregationStore,
	categoryAggregationStore CategoryAggregationStore,
) *Storage {
	return &Storage{
		CategoryStore:            categoryStore,
		RuleStore:                ruleStore,
		AppAggregationStore:      appAggregationStore,
		CategoryAggregationStore: categoryAggregationStore,
	}
}
