package config

import (
	"encoding/json"
	"log"
	"mediator/pkg/logger"

	"os"
)

var (
	Params  Configuration
	Logging logger.Logger
)

func init() {
	var err error
	Params, err = ReadConcfig()
	if err != nil {
		log.Println("unknown configuration. Check the config file")
		panic(err)
	}
}

type Configuration struct {
	DBLogin      string
	DBPass       string
	DBHost       string
	DBPort       int
	RPCHost      string
	RPCPort      int
	DBName       string
	MediatorName string
}

func ReadConcfig() (Configuration, error) {
	config := Configuration{}
	var err error
	config_file := "conf.json"
	_, err = os.Stat(config_file)
	if err == nil {
		file, _ := os.Open(config_file)
		defer file.Close()
		decoder := json.NewDecoder(file)
		err = decoder.Decode(&config)
	}
	return config, err
}
