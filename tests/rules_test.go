package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

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

func postRequest(body, route_name string) (*http.Response, error) {
	rules_uri := fmt.Sprintf("%s/%s", uri, route_name)
	request, err := http.NewRequest("POST", rules_uri, bytes.NewBuffer([]byte(body)))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	return response, err

}

func TestPostRules(t *testing.T) {

	rulesBody := rulesContainer()
	for rule_name, body := range rulesBody {
		response, err := postRequest(body, "rules")
		if response == nil {
			t.Error(err)
		}

		if rule_name == "rule1" && response.StatusCode != 200 {
			t.Error(errors.New("rule1 incorrect actions"))
		}

		if rule_name == "rule2" && response.StatusCode == 200 {
			t.Error(errors.New("rule2 incorrect actions"))
		}
		if rule_name == "rule3" && response.StatusCode == 200 {
			t.Error(errors.New("rule3 incorrect actions"))
		}
		if rule_name == "rule4" && response.StatusCode != 200 {
			t.Error(errors.New("rule4 incorrect actions"))
		}

	}

}

func TestGetRules(t *testing.T) {
	rules_uri := fmt.Sprintf("%s/rules", uri)
	resp, err := http.Get(rules_uri)
	if err != nil {
		t.Error(err)
		return
	}
	ruledata := new(ModelRules)
	if resp.StatusCode != 200 {
		t.Errorf("Get rules bad status code: %d", resp.StatusCode)
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &ruledata)
	if err != nil {
		t.Errorf("parsing rules json response: %v", err)
		return
	}

}

func TestGetRule(t *testing.T) {

	rulesBody := rulesContainer()
	postRequest(rulesBody["rule13"], "rules")
	rules_uri := fmt.Sprintf("%s/rules/22.22.22", uri)
	resp, err := http.Get(rules_uri)
	if err != nil {
		t.Error(err)
		return
	}

	if resp.StatusCode != 200 {
		t.Errorf("Get rule with prefix bad status code: %d", resp.StatusCode)
		return
	}
	rules_uri = fmt.Sprintf("%s/rules/1.1.1", uri)
	resp, err = http.Get(rules_uri)
	if err != nil {
		t.Error(err)
		return
	}
	if resp.StatusCode != 404 {
		t.Errorf("Find prefix 1.1.1")
		return
	}

	ruledata := new(ModelRules)
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &ruledata)
	if err != nil {
		t.Errorf("parsing rule with prefix json response: %v", err)
		return
	}

}
