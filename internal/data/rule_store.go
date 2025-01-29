package data

type RuleStore interface {
	GetRules() ([]CategoryRule, error)
	GetRulesByCategory(categoryId CategoryId) ([]CategoryRule, error)
	GetRulesByApp(appName string) ([]CategoryRule, error)
	AddRule(rule CategoryRule) error
	DeleteRule(rule CategoryRule) error
}

type RuleStore_InMemory_Impl struct {
	data []CategoryRule
}

func NewRuleStore_InMemory_Impl() *RuleStore_InMemory_Impl {
	return &RuleStore_InMemory_Impl{
		data: make([]CategoryRule, 0),
	}
}

func (r *RuleStore_InMemory_Impl) GetRules() ([]CategoryRule, error) {
	return r.data, nil
}

func (r *RuleStore_InMemory_Impl) GetRulesByCategory(categoryId CategoryId) ([]CategoryRule, error) {
	panic("implement me")
}

func (r *RuleStore_InMemory_Impl) GetRulesByApp(appName string) ([]CategoryRule, error) {
	res := make([]CategoryRule, 0)
	for _, rule := range r.data {
		if rule.AppName == appName {
			res = append(res, rule)
		}
	}
	return res, nil
}

func (r *RuleStore_InMemory_Impl) AddRule(rule CategoryRule) error {
	r.data = append(r.data, rule)
	return nil
}

func (r *RuleStore_InMemory_Impl) DeleteRule(rule CategoryRule) error {
	panic("implement me")
}
