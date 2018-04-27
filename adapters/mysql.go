package adapters

import (
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql" // mysql
)

const createTableMySQL = "CREATE TABLE `_migrations` (`name` varchar(255) NOT NULL DEFAULT '', `migrated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, PRIMARY KEY (`name`)) ENGINE=InnoDB DEFAULT CHARSET=utf8;"

// MySQL - DB adapter
type MySQL struct {
	config     Config
	conn       *sql.DB
	migrations migrations
	processed  map[string]Migration
	SQLAdapter
}

// NewMySQL - MySQL adapter constructor
func NewMySQL(config Config, m migrations) *MySQL {
	db := &MySQL{migrations: m}
	db.config = db.prepareConfig(config)
	db.SQLAdapter = SQLAdapter{config: db.config}
	return db
}

func (db *MySQL) prepareConfig(config Config) Config {
	if config.RTimeout == 0 {
		config.RTimeout = 1
	}
	if config.RAttempts == 0 {
		config.RAttempts = 1
	}
	if config.User == "" {
		config.User = "root"
	}
	if config.Host == "" {
		config.Host = "localhost"
	}
	if config.Port == 0 {
		config.Port = 3306
	}
	return config
}

func (db *MySQL) getCreateTableQuery() string {
	return createTableMySQL
}

// getConnectionString - default MySQL
func (db *MySQL) getConnectionString(database string) string {
	return fmt.Sprintf("%v:%v@tcp(%v:%v)/%v",
		db.config.User,
		db.config.Password,
		db.config.Host,
		db.config.Port,
		database,
	)
}

// Migrate - starting migration
func (db *MySQL) Migrate() (err error) {
	if err = db.connect(db.config.Database); err != nil {
		return
	}
	if err = db.loadProcessed(); err != nil {
		return
	}

	var keys []string
	for key := range db.migrations {
		if _, ok := db.processed[key]; ok == false {
			keys = append(keys, key)
		}
	}
	sort.Strings(keys)
	for _, key := range keys {
		if err = db.migrate(db.migrations[key]); err != nil {
			return
		}
		if err = db.save(key); err != nil {
			return
		}
	}
	return
}

// check - checking for migrations
func (db *MySQL) loadProcessed() (err error) {
	var SQL string
	migrations := make(map[string]Migration)

	SQL = `SELECT name, migrated_at FROM ` + db.config.Database + `.` + table + `;`
	rows, err := db.conn.Query(SQL)
	if err != nil {
		fmt.Println("Error loading migrations from DB: ", err)
		fmt.Println("Will try to create")
		if db.checkError(err) == 2 {
			if err = db.createTable(); err != nil {
				return
			}
			return db.loadProcessed()
		}
		return
	}
	defer rows.Close()

	for rows.Next() {
		var m Migration
		var migratedAt string
		if err = rows.Scan(&m.Name, &migratedAt); err != nil {
			return
		}
		m.MigratedAt, _ = time.Parse("2006-01-02 15:04:05", migratedAt)
		migrations[m.Name] = m
	}
	db.processed = migrations
	return
}

func (db *MySQL) createTable() (err error) {
	fmt.Println("Creating table: ", table)
	_, err = db.conn.Exec(db.getCreateTableQuery())
	if err != nil {
		fmt.Println("Error creating table: ", err)
		return
	}
	return
}

func (db *MySQL) createDB() (err error) {
	fmt.Println("Creating database: ", db.config.Database)
	db.connect("")
	SQL := `CREATE DATABASE ` + db.config.Database + `;`
	_, err = db.conn.Exec(SQL)
	if err != nil {
		fmt.Println("Error creating DB: ", err)
		return
	}
	if err = db.conn.Close(); err != nil {
		fmt.Println("Error closing connection: ", err)
		return
	}
	return db.connect(db.config.Database)
}

func (db *MySQL) save(key string) (err error) {
	migratedAt := time.Now().Format("2006-01-02T15:04:05")
	SQL := `INSERT INTO ` + table + ` (name, migrated_at) VALUES ('` + key + `', '` + migratedAt + `');`
	_, err = db.conn.Exec(SQL)
	return
}

func (db *MySQL) migrate(data string) (err error) {
	var queries []string
	if queries, err = parser.ParseFromString(data); err != nil {
		return
	}
	tx, _ := db.conn.Begin()
	defer tx.Commit()
	for _, query := range queries {
		if _, err = tx.Exec(query); err != nil {
			fmt.Println("[ERROR] Error: ", err)
			fmt.Println("[ERROR] SQL: ", query)
			tx.Rollback()
			return
		}
	}
	return
}

func (db *MySQL) connect(database string) (err error) {
	if db.conn != nil && db.conn.Stats().OpenConnections > 0 {
		return
	}
	config := db.config
	if config.Database == "" {
		return errors.New("Database is empty")
	}

	connectionInfo := db.getConnectionString(database)
	conn, err := db.open(config.Driver, connectionInfo)
	if err != nil {
		return fmt.Errorf("Connection error %s", err)
	}
	if conn == nil {
		return fmt.Errorf("Connection is nil")
	}
	if err = conn.Ping(); err != nil {
		if db.checkError(err) == 1 {
			return db.createDB()
		}

		if db.config.RAttempts > reconnectAttempts {
			reconnectAttempts++
			fmt.Printf("Failed connection: %v\n", err)
			fmt.Printf("Reconnecting %d of %d ...\n", reconnectAttempts, db.config.RAttempts)
			time.Sleep(time.Duration(db.config.RTimeout) * time.Second)
			return db.connect(database)
		}
		return fmt.Errorf("Failed connection: %v", err)
	}
	db.conn = conn
	return
}
func (db *MySQL) checkError(err error) (code int) {
	if strings.Contains(err.Error(), fmt.Sprintf("Error 1049: Unknown database '%s'", db.config.Database)) {
		return 1
	}
	if strings.Contains(err.Error(), fmt.Sprintf("Error 1146: Table '%s.%s' doesn't exist", db.config.Database, table)) {
		return 2
	}
	return 0
}

func (db *MySQL) open(driver, connection string) (*sql.DB, error) {
	return sql.Open(driver, connection)
}
