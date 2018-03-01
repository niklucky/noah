package main

import (
	"encoding/json"

	lib "github.com/niklucky/go-lib"
	"github.com/niklucky/noah/adapters"
)

func parseDBConfig(configFileName string) (config adapters.Config, err error) {
	fileData, err := lib.ReadFile(configFileName)
	if err != nil {
		return
	}
	if err = json.Unmarshal(fileData, &config); err != nil {
		return
	}
	if config.User == "" {
		config.User = "root"
	}
	if config.Host == "" {
		config.Host = "localhost"
	}
	return config, err
}
