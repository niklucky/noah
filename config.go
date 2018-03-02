package main

import (
	lib "github.com/niklucky/go-lib"
	"github.com/niklucky/noah/adapters"
	"gopkg.in/yaml.v2"
)

func parseDBConfig(configFileName string) (config adapters.Config, err error) {
	fileData, err := lib.ReadFile(configFileName)
	if err != nil {
		return
	}
	if err = yaml.Unmarshal(fileData, &config); err != nil {
		return
	}
	return config, err
}
