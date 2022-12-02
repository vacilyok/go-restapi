package main

import (
	"flag"
	"log"
	"net/http"

	devicehandlers "mediator/internal/adapters/api/devices"
	metricshandler "mediator/internal/adapters/api/metrics"
	ruleshandlers "mediator/internal/adapters/api/rules"
	"mediator/internal/config"
	"mediator/internal/domain/devices"
	"mediator/internal/domain/rules"
	"mediator/internal/domain/watcher"
	mysqldb "mediator/pkg/database/mysql"
	"mediator/pkg/logger"

	mux "gitlab.ddos-guard.net/dma/gorilla"
	handle "gitlab.ddos-guard.net/dma/gorilla-handlers"
)

func start(router *mux.Router) {
	addr := flag.String("addr", ":4000", "Web server network address")
	flag.Parse()
	headersOk := handle.AllowedHeaders([]string{"Content-Type"})
	methodsOk := handle.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})
	originsOK := handle.AllowedOrigins([]string{"*"})
	config.Mysqllog.Info("Start mediator service")
	http.ListenAndServe(*addr, handle.CORS(originsOK, headersOk, methodsOk)(router))

}

func main() {
	dbConn, err := mysqldb.DBconnect()
	isConnection := false
	if err != nil {
		log.Println("ERROR: Not connect to database")
	} else {
		mysqldb.InitDB(dbConn)
		isConnection = true
	}
	defer dbConn.Close()

	config.Mysqllog = logger.NewLogger(dbConn, isConnection)
	router := mux.NewRouter()
	devStoreage := devices.NewDevStorage(dbConn, isConnection)
	devService := devices.NewService(devStoreage)
	devHandlers := devicehandlers.NewDevHandlers(devService)
	devHandlers.Register(router)

	ruleStorage := rules.NewRuleStorage(dbConn, isConnection)
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
