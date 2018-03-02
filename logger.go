package main

import (
	"log"

	"github.com/fatih/color"
)

const (
	infoLoadingConfig     = "Loading config from file"
	infoDatabaseStart     = "Starting %s migration"
	infoConnectingToDB    = "Connecting to DB %s"
	infoConnectedToDB     = "Connected to DB %s"
	infoLoadingMigrations = "Loading migrations from: %s"
)

func info(str string) {
	log.Printf("[%s]: %s\n", color.HiCyanString("INFO"), str)
}

func fatal(str string) {
	log.Fatalf("[%s]: %s\n", color.HiRedString("FATAL"), str)
}
