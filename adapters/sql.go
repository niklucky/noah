package adapters

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql" // mysql
)

// SQLAdapter - DB adapter
type SQLAdapter struct {
	config     Config
	conn       *sql.DB
	migrations migrations
	processed  map[string]Migration
}
