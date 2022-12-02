package main

import (
	"context"
	"flag"
	"net/http"

	devicehandlers "mediator/internal/adapters/api/devices"
	metricshandler "mediator/internal/adapters/api/metrics"
	ruleshandlers "mediator/internal/adapters/api/rules"
	"mediator/internal/config"
	"mediator/internal/domain/devices"
	"mediator/internal/domain/rules"
	"mediator/internal/domain/watcher"
	"mediator/pkg/database/pg"
	"mediator/pkg/logger"

	mux "github.com/gorilla"
	handle "github.com/gorilla-handlers"
)

func start(router *mux.Router) {
	addr := flag.String("addr", ":4000", "Web server network address")
	flag.Parse()
	headersOk := handle.AllowedHeaders([]string{"Content-Type"})
	methodsOk := handle.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})
	originsOK := handle.AllowedOrigins([]string{"*"})
	config.Logging.Info("Start mediator service")
	if err := http.ListenAndServe(*addr, handle.CORS(originsOK, headersOk, methodsOk)(router)); err != nil {
		config.Logging.Error(err.Error())
	}
}

func main() {
	config.Logging = logger.NewLogger()
	isConnection := true
	pgxConn, err := pg.DBInit()
	if err != nil {
		config.Logging.Error(err.Error())
		isConnection = false
	}
	defer pgxConn.Close(context.Background())
	router := mux.NewRouter()
	devStorage := devices.NewDevStorage(pgxConn)
	devService := devices.NewService(devStorage)
	devHandlers := devicehandlers.NewDevHandlers(devService)
	devHandlers.Register(router)

	ruleStorage := rules.NewRuleStorage(pgxConn)
	ruleService := rules.NewService(ruleStorage)
	ruleHandlers := ruleshandlers.NewRulesHandlers(ruleService)
	ruleHandlers.Register(router)

	MetricsHandler := metricshandler.NewMetricsHandlers()
	MetricsHandler.Register(router)
	if isConnection {
		watch := watcher.NewWatcher(devService, ruleService)
		watch.StartWatch()
	}
	start(router)
}
