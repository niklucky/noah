package main

import "flag"

type settings struct {
	configFile string
	dir        string
	mode       string
}

func parseFlags() settings {
	var configFile string
	var dumpDir string
	var mode string

	var (
		modeDev,
		modeProduction,
		modeTesting bool
	)

	flag.StringVar(&configFile, "config", "./config.yml", "Pass config file where DB credentials are")
	flag.StringVar(&dumpDir, "dir", "./migrations", "Directory where Dump files are located")
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
	return settings{
		dir:        dumpDir,
		configFile: configFile,
		mode:       mode,
	}
}
