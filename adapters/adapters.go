package adapters

import (
	"errors"
	"fmt"
	"time"

	sqlparser "github.com/PGV65/sql-parser"
)

const table = "_migrations"

var parser = sqlparser.Parser{}
var reconnectAttempts = 0

// Config - database config
type Config struct {
	Driver    string `json:"driver"`
	Host      string `json:"host"`
	Port      int    `json:"port"`
	User      string `json:"user"`
	Password  string `json:"password"`
	Database  string `json:"database"`
	SSLmode   string `json:"ssl_mode" yaml:"ssl_mode"`
	RTimeout  int    `json:"reconnect_timeout" yaml:"reconnect_timeout"`
	RAttempts int    `json:"reconnect_attempts" yaml:"reconnect_attempts"`
}

type migrations map[string]string

// Migration - enity for DB storing. Name is a fileName from dump
type Migration struct {
	Name       string
	MigratedAt time.Time
}

// Adapter — DB adapter interface
type Adapter interface {
	Migrate() error
}

// New - adapter constructor
func New(config Config, migrations map[string]string) (Adapter, error) {

	if config.Driver == "mysql" {
		return NewMySQL(config, migrations), nil
	}
	if config.Driver == "postgres" {
		return NewPostgres(config, migrations), nil
	}
	return nil, errors.New("Driver not found")
}

func info(values ...interface{}) {
	fmt.Println(values...)
}
