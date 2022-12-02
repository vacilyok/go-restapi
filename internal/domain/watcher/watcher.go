package watcher

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"mediator/internal/config"
	"mediator/internal/domain/devices"
	"mediator/internal/domain/rules"
	"net/http"
	"strconv"
	"time"
)

type watcher struct {
	devservice devices.DevService
	rulservice rules.RuleService
}

// create new rules handler
func NewWatcher(dservice devices.DevService, rservice rules.RuleService) WatchService {
	return &watcher{
		devservice: dservice,
		rulservice: rservice,
	}
}

// *******************************************************************************
func (w *watcher) StartWatch() {
	var startMediatorTime float64 = 0
	ticker := time.NewTicker(time.Second)
	go func() {
		for range ticker.C {
			uptimeSystem, err := w.getUptime()
			if err != nil {
				config.Logging.Error(err.Error())
				continue
			}
			if startMediatorTime == 0 {
				startMediatorTime = 1
				config.Logging.Info("Launched initialization process at startup")
				err = w.devservice.InitDevices()
				if err != nil {
					config.Logging.Error(err.Error())
					return
				}
				err = w.rulservice.InitRules()
				if err != nil {
					config.Logging.Error(err.Error())
					return
				}
				continue
			}
			isRestart := w.isRestart(&startMediatorTime, uptimeSystem)
			if isRestart {
				config.Logging.Info("Service restart detected. Сonfiguration process started ...")
				err = w.devservice.InitDevices()
				if err != nil {
					config.Logging.Error(err.Error())
					return
				}
				w.rulservice.InitRules()
				if err != nil {
					config.Logging.Error(err.Error())
					return
				}

			}

		}
	}()
}

// *******************************************************************************
func (w *watcher) isRestart(startMediatorTime *float64, uptime float64) bool {

	if uptime < *startMediatorTime {
		*startMediatorTime = uptime
		return true
	}
	if *startMediatorTime <= uptime {
		*startMediatorTime = uptime
		return false
	}
	*startMediatorTime = uptime
	return true
}

// *******************************************************************************
func (w *watcher) getUptime() (float64, error) {
	url := "http://" + config.Params.RPCHost + ":" + strconv.Itoa(config.Params.RPCPort) + "/system"
	response, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	body, _ := ioutil.ReadAll(response.Body)
	systemMetrics := make(map[string]interface{})
	err = json.Unmarshal(body, &systemMetrics)
	if err != nil {
		return 0, err
	}
	if uptime, ok := systemMetrics["uptime"]; ok {
		return uptime.(float64), nil
	}

	return 0, errors.New(" no key uptime in system")

}
