package main

import (
	"log"
	"os"
	"strconv"

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
	parseEnv(&config)
	return config, err
}

func parseEnv(config *adapters.Config) {
	if env, ok := os.LookupEnv("DB_DRIVER"); ok {
		config.Driver = env
	}
	if env, ok := os.LookupEnv("DB_HOST"); ok {
		config.Host = env
	}
	if env, ok := os.LookupEnv("DB_PORT"); ok {
		config.Port = envToInt(env)
	}
	if env, ok := os.LookupEnv("DB_USER"); ok {
		config.User = env
	}
	if env, ok := os.LookupEnv("DB_PASSWORD"); ok {
		config.Password = env
	}
	if env, ok := os.LookupEnv("DB_DATABASE"); ok {
		config.Database = env
	}
	if env, ok := os.LookupEnv("DB_SSLMODE"); ok {
		config.SSLmode = env
	}
	if env, ok := os.LookupEnv("DB_RTIMEOUT"); ok {
		config.RTimeout = envToInt(env)
	}
	if env, ok := os.LookupEnv("DB_RATTEMPTS"); ok {
		config.RAttempts = envToInt(env)
	}
}

func envToInt(s string) (i int) {
	var err error
	i, err = strconv.Atoi(s)
	if err != nil {
		log.Println("Error in DB_PORT: ", err)
	}
	return
}
