package rules

import (
	"net/http"
)

type RuleService interface {
	CreateNewRule(w http.ResponseWriter, r *http.Request) (*http.Response, error)
	GetRules(r *http.Request, prefix string) (*http.Response, error)
	InitRules() error
}
