package ruleshandlers

import (
	"io/ioutil"
	"mediator/internal/adapters/api"
	"mediator/internal/config"
	"mediator/internal/domain/rules"
	"net/http"

	mux "gitlab.ddos-guard.net/dma/gorilla"
)

const (
	rulesURL = "/rules"
	ruleURL  = "/rules/{prefix}"
)

type rulesHandler struct {
	service rules.RuleService
}

// create new rules handler
func NewRulesHandlers(service rules.RuleService) api.Handlers {
	return &rulesHandler{
		service: service,
	}
}

func (rh *rulesHandler) viewResponse(writer http.ResponseWriter, resp *http.Response) {
	contentType := resp.Header.Get("Content-type")
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		config.Mysqllog.Error(err.Error())
		rh.errorResponse(writer, err, 500)
		return
	}
	writer.Header().Set("Content-Type", contentType)
	writer.WriteHeader(resp.StatusCode)
	writer.Write(respBody)
}

// http error message response
func (rh *rulesHandler) errorResponse(w http.ResponseWriter, err error, code int) {
	w.WriteHeader(code)
	w.Write([]byte(err.Error()))
}

func (rh *rulesHandler) Register(router *mux.Router) {
	router.HandleFunc(rulesURL, rh.Create).Methods("POST")
	router.HandleFunc(rulesURL, rh.GetRules).Methods("GET")
	router.HandleFunc(ruleURL, rh.GetRule).Methods("GET")
}

// handler for create new rules
func (rh *rulesHandler) Create(w http.ResponseWriter, r *http.Request) {
	response, err := rh.service.CreateNewRule(w, r)
	if err != nil {
		rh.errorResponse(w, err, 500)
		return
	}
	rh.viewResponse(w, response)
}

func (rh *rulesHandler) GetRule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if _, ok := vars["prefix"]; ok {
		prefix := vars["prefix"]
		response, err := rh.service.GetRules(r, prefix)
		if err != nil {
			rh.errorResponse(w, err, 500)
			return
		}
		rh.viewResponse(w, response)
	}

}

func (rh *rulesHandler) GetRules(w http.ResponseWriter, r *http.Request) {

	response, err := rh.service.GetRules(r, "")
	if err != nil {
		rh.errorResponse(w, err, 500)
		return
	}
	rh.viewResponse(w, response)
}
