package metricshandler

import (
	"fmt"
	"io/ioutil"
	"mediator/internal/adapters/api"
	"mediator/internal/config"
	"net/http"

	mux "gitlab.ddos-guard.net/dma/gorilla"
)

type metricHandler struct {
}

// new handler for metrics
func NewMetricsHandlers() api.Handlers {
	return &metricHandler{}
}

// *******************************************************************************************
// register metrics handlers
func (mh *metricHandler) Register(router *mux.Router) {
	router.HandleFunc("/metrics", mh.GetMetrics).Methods("GET")
	router.HandleFunc("/system", mh.GetSystem).Methods("GET")
}

// *******************************************************************************************
// send GET request rpc server for receiving all metrics
func (mh *metricHandler) GetMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	url := fmt.Sprintf("http://%s:%d/metrics", config.Params.RPCHost, config.Params.RPCPort)
	resp, err := http.Get(url)
	if err != nil {
		w.WriteHeader(503)
		w.Write([]byte(err.Error()))
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		w.WriteHeader(503)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(body)
}

// *******************************************************************************************
func (mh *metricHandler) GetSystem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	url := fmt.Sprintf("http://%s:%d/system", config.Params.RPCHost, config.Params.RPCPort)
	resp, err := http.Get(url)
	if err != nil {
		w.WriteHeader(503)
		w.Write([]byte(err.Error()))
		return
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		w.WriteHeader(503)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(resp.StatusCode)
	w.Write(respBody)
}
