package rulestorage

type ModelRules struct {
	Rules   []Rule       `json:"rules"`
	Buckets []BucketItem `json:"buckets"`
	Lists   []ListItem   `json:"lists"`
}

type Rule struct {
	Prefix      string           `json:"prefix"`
	Countermeas []Countermeasura `json:"countermeasures"`
}

type Countermeasura struct {
	Matches []MatchesItem `json:"matches"`
	Action  ActionItem    `json:"action"`
}

type ActionItem struct {
	Name    string      `json:"name"`
	Options interface{} `json:"options"`
}

type MatchesItem struct {
	Name    string                 `json:"name"`
	Options map[string]interface{} `json:"options"`
}

type BucketItem struct {
	Name      string `json:"name"`
	Limit_bps int    `json:"limit_bps"`
	Limit_pps int    `json:"limit_pps"`
}

type ListItem struct {
	Name  string   `json:"name"`
	Items []string `json:"items"`
}
