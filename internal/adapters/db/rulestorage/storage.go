package rulestorage

type RuleStorage interface {
	SaveRulesToDB(ruleList map[string]interface{}) error
	GetRulesFromDB() (string, error)
}
