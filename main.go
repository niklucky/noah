package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/niklucky/noah/adapters"
)

var flags settings

func main() {
	flags = parseFlags()

	intro()

	info(infoLoadingConfig)

	if flags.isServer {
		fmt.Println("Running in HTTP-server mode")
		startServer(flags.port)
		return
	}
	config, err := parseDBConfig(flags.configFile)
	if err != nil {
		fatal(fmt.Sprintln("Error parsing config: ", err))
	}

	info(fmt.Sprintf(infoLoadingMigrations, color.HiYellowString(flags.dir)))
	migrations, err := loadMigrations(flags.dir, flags.mode)
	if err != nil {
		fatal(fmt.Sprintln("Error loading migration data: ", err))
	}
	info(fmt.Sprintf(infoDatabaseStart, config.Database))
	info(fmt.Sprintf(infoConnectingToDB, color.HiYellowString(config.Database)))

	adapter, err := adapters.New(config, migrations)
	if err != nil {
		fatal(fmt.Sprintf("Error creating adapter: %v", err))
	}
	info(fmt.Sprintf(infoConnectedToDB, color.HiYellowString(config.Database)))

	if err = adapter.Migrate(); err != nil {
		fatal(fmt.Sprintf("Error in migration:  %v", err))
	}
	info("Success! Migrations finished.")
}

func intro() {
	fmt.Println(color.HiCyanString("===================================================="))
	fmt.Printf("\n\t\t%s: Starting migrations\n\n", color.HiGreenString("Noah 0.1.0"))
	fmt.Println(color.HiCyanString("===================================================="))
	fmt.Printf("\n[INFO] Config file: \t%s\n", flags.configFile)
	fmt.Printf("[INFO] Dump dir: \t\t%s\n", flags.dir)
	fmt.Printf("[INFO] Mode: \t\t%s\n\n", color.HiMagentaString(flags.mode))
	fmt.Println(color.HiCyanString("----------------------------------------------------"))
}
