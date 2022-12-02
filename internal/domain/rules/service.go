package rules

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mediator/internal/adapters/db/rulestorage"
	"mediator/internal/config"
	"net/http"
	"strconv"
)

const rulesURL string = "/rules"

type service struct {
	storage rulestorage.RuleStorage
}

// Create new service for rules
func NewService(storage rulestorage.RuleStorage) RuleService {
	return &service{
		storage: storage,
	}
}

func (s *service) InitRules() error {
	config.Logging.Info("Start init rules")
	rulesBody, err := s.storage.GetRulesFromDB()
	if err != nil {
		config.Logging.Warning("No init rules ")
		return err
	}
	if len(rulesBody) == 0 {
		config.Logging.Info("Rules empty. Initialization empty rules")
		rulesBody = `{"rules": []}`
	}

	url := "http://" + config.Params.RPCHost + ":" + strconv.Itoa(config.Params.RPCPort) + rulesURL
	config.Logging.LogRest(rulesBody, rulesURL, "POST")
	request, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(rulesBody)))
	if err != nil {
		config.Logging.Error("post request to rpc server is fail")
		return err
	}
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode == 200 {
		config.Logging.Info("Successful initialization of rules")
	}
	return nil
}

// *************************************************************************************************
// Service create new rule
func (s *service) CreateNewRule(w http.ResponseWriter, r *http.Request) (*http.Response, error) {
	ruleBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		config.Logging.Error("Error parse rule body request")
		return nil, err
	}
	var rules map[string]interface{}
	if err = json.Unmarshal(ruleBody, &rules); err != nil {
		config.Logging.Error(err.Error())
		return nil, err
	}
	response, err := s.postRuleRequestInterface(&rules)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		err_msg := make(map[string]string)
		body, _ := ioutil.ReadAll(response.Body)
		json.Unmarshal(body, &err_msg)
		config.Logging.Error(err_msg["error"])
		return nil, errors.New(err_msg["error"])
	}
	err = s.storage.SaveRulesToDB(rules)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// *************************************************************************************************
func (s *service) postRuleRequestInterface(rules *map[string]interface{}) (*http.Response, error) {
	// post request to RPC server
	jsonString, _ := json.Marshal(rules)
	url := "http://" + config.Params.RPCHost + ":" + strconv.Itoa(config.Params.RPCPort) + rulesURL
	config.Logging.LogRest(string(jsonString), rulesURL, "POST")
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonString))
	if err != nil {
		config.Logging.Error("post request to rpc server is fail")
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	client := &http.Client{}
	response, err := client.Do(request)
	return response, err
}

// *************************************************************************************************
func (s *service) GetRules(r *http.Request, prefix string) (*http.Response, error) {
	url := fmt.Sprintf("http://%s:%d%s/%s", config.Params.RPCHost, config.Params.RPCPort, rulesURL, prefix)
	resp, err := http.Get(url)
	if err != nil {
		msg := "get rules from RPC server is fail"
		config.Logging.Error(errors.New(msg).Error())
		return nil, err
	}
	return resp, nil
}
