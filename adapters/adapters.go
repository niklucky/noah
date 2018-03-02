package adapters

import (
	"errors"
	"time"
)

// Config - database config
type Config struct {
	Driver   string `json:"driver"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

type migrations map[string]string

// Migration - enity for DB storing. Name is a fileName from dump
type Migration struct {
	Name       string
	MigratedAt time.Time
}

// Adapter — DB adapter interface
type Adapter interface {
	Check() error
	Migrate() error
}

// New - adapter constructor
func New(config Config, migrations map[string]string) (Adapter, error) {

	if config.Driver == "mysql" {
		return NewMySQL(config, migrations), nil
	}
	return nil, errors.New("Driver not found")
}
