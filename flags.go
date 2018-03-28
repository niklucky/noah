package main

import "flag"

type settings struct {
	configFile string
	dir        string
	mode       string
	isServer   bool
	port       int
}

func parseFlags() settings {
	var configFile string
	var dumpDir string
	var mode string
	var port int

	var isServer bool

	flag.StringVar(&configFile, "config", "./config.yml", "Pass config file where DB credentials are")
	flag.StringVar(&dumpDir, "dir", "./migrations", "Directory where Dump files are located")

	flag.BoolVar(&isServer, "s", false, "Running in HTTP-server mode")
	flag.IntVar(&port, "port", 12000, "Port for HTTP-server")

	flag.StringVar(&mode, "mode", "production", "Migration mode prefix")

	flag.Parse()

	return settings{
		dir:        dumpDir,
		configFile: configFile,
		isServer:   isServer,
		port:       port,
		mode:       mode,
	}
}
