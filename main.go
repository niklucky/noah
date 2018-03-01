package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/niklucky/noah/adapters"
)

func main() {
	var configFile string
	var dumpDir string
	var mode string
	var (
		modeDev,
		modeProduction,
		modeTesting bool
	)

	flag.StringVar(&configFile, "config", ".", "Pass config file where DB credentials are")
	flag.StringVar(&dumpDir, "dir", ".", "Directory where Dump files are located")
	flag.BoolVar(&modeDev, "D", true, "Development staging migration")
	flag.BoolVar(&modeProduction, "P", false, "Production staging migration")
	flag.BoolVar(&modeTesting, "T", false, "Testing staging migration")

	flag.Parse()

	if modeProduction {
		mode = "prod"
	}
	if modeDev {
		mode = "dev"
	}
	if modeTesting {
		mode = "test"
	}
	fmt.Println("Config file: ", configFile)
	fmt.Println("Dump dir: ", dumpDir)
	fmt.Println("Mode: ", mode)
	config, err := parseDBConfig(configFile)
	if err != nil {
		log.Fatalln("Error parsing config: ", err)
	}
	fmt.Println("Config: ", config)
	adapter, err := adapters.New(config)
	if err != nil {
		log.Fatalln("Error creating adapter: ", err)
	}
	migrations, err := loadMigrations(dumpDir, mode)
	if err != nil {
		log.Fatalln("Error loading data: ", err)
	}
	fmt.Println("Migrations: ", len(migrations))
	adapter.AddMigrations(migrations)
	var processed map[string]adapters.Migration
	if processed, err = adapter.Check(); err != nil {
		fmt.Println("Error in check migration: ", err)
	}
	fmt.Println("Processed: ", processed)
	if err = adapter.Migrate(); err != nil {
		fmt.Println("Error in migration: ", err)
	}
	fmt.Println("Success! Migrations finished.")
}
