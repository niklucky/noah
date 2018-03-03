package adapters

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	sqlparser "github.com/PGV65/sql-parser"
	_ "github.com/go-sql-driver/mysql" // mysql
)

const table = "_migrations"
const createTableSQL = "CREATE TABLE `_migrations` (`name` varchar(255) NOT NULL DEFAULT '', `migrated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, PRIMARY KEY (`name`)) ENGINE=InnoDB DEFAULT CHARSET=utf8;"

var parser = sqlparser.Parser{}
var reconnectAttempts = 0

// MySQL - DB adapter
type MySQL struct {
	config Config
	conn   *sql.DB
	SQLAdapter
	migrations migrations
	processed  map[string]Migration
}

// NewMySQL - MySQL adapter constructor
func NewMySQL(config Config, m migrations) *MySQL {
	fmt.Printf("Config: %+v\n\n", config)
	if config.RTimeout == 0 {
		config.RTimeout = 1
	}
	if config.RAttempts == 0 {
		config.RAttempts = 1
	}
	return &MySQL{
		config:     config,
		migrations: m,
	}
}

// check - checking for migrations
func (db *MySQL) check() (err error) {
	if err = db.connect(); err != nil {
		return
	}

	var SQL string
	migrations := make(map[string]Migration)

	SQL = `SELECT name, migrated_at FROM ` + db.config.Database + `.` + table + `;`
	rows, err := db.exec(SQL)
	if err != nil {
		fmt.Println("Error loading migrations from DB: ", err)
		if strings.Contains(err.Error(), "Error 1049") {
			if err = db.createDB(); err != nil {
				return
			}
			return db.check()
		} else if strings.Contains(err.Error(), "Error 1146") {
			if err = db.createDB(); err != nil {
				return
			}
			if err = db.createTable(); err != nil {
				return
			}
			return db.check()
		} else {
			return
		}
	} else {
		defer rows.Close()
	}

	err = nil

	SQL = `USE ` + db.config.Database + `;`
	_, err = db.conn.Exec(SQL)
	if err != nil {
		fmt.Println("Error while selecting DB: ", err)
		return
	}
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
	_, err = db.conn.Exec(`USE ` + db.config.Database + `;`)
	if err != nil {
		fmt.Println("Error creating table: ", err)
		return
	}
	_, err = db.conn.Exec(createTableSQL)
	if err != nil {
		fmt.Println("Error creating table: ", err)
		return
	}
	return
}

func (db *MySQL) createDB() (err error) {
	SQL := `CREATE DATABASE ` + db.config.Database + `;`
	fmt.Println(SQL)

	_, err = db.conn.Exec(SQL)
	if err != nil {
		fmt.Println("Error creating DB: ", err)
		return
	}
	return
}

// Migrate - starting migration
func (db *MySQL) Migrate() (err error) {
	if err = db.connect(); err != nil {
		return
	}
	if err = db.check(); err != nil {
		return
	}

	var keys []string
	for key := range db.migrations {
		if _, ok := db.processed[key]; ok == false {
			keys = append(keys, key)
		}
	}
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

func (db *MySQL) save(key string) (err error) {
	SQL := `INSERT INTO ` + db.config.Database + `.` + table + ` (name) VALUES ('` + key + `');`
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
	tx.Exec("USE " + db.config.Database)
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

/*
Exec - executing SQL-query and returning *Rows
*/
func (db *MySQL) exec(SQL string) (*sql.Rows, error) {
	if err := db.checkConnection(); err != nil {
		return nil, err
	}
	return db.conn.Query(SQL)
}

/*
query - preparing query into Statement and executing SQL-query and returning *Rows
*/
func (db *MySQL) query(SQL string, values []interface{}) (*sql.Rows, error) {
	if err := db.checkConnection(); err != nil {
		return nil, err
	}
	return db.conn.Query(SQL, values...)
}

func (db *MySQL) connect() error {
	config := db.config
	if config.User == "" {
		config.User = "root"
	}
	if config.Host == "" {
		config.Host = "localhost"
	}
	if config.Port == 0 {
		config.Port = 3306
	}
	if config.Database == "" {
		return errors.New("Database is empty")
	}
	connectionInfo := fmt.Sprintf("%v:%v@tcp(%v:%v)/",
		config.User,
		config.Password,
		config.Host,
		config.Port,
	)
	conn, err := sql.Open("mysql", connectionInfo)
	if err != nil {
		fmt.Println("MySQL connection error", err)
		return err
	}
	if conn == nil {
		fmt.Println("Connection to MySQL is nil")
	}
	if err = conn.Ping(); err != nil {
		if db.config.RAttempts > reconnectAttempts {
			reconnectAttempts++
			fmt.Printf("Failed connection: %v\n", err)
			fmt.Printf("Reconnecting %d of %d ...\n", reconnectAttempts, db.config.RAttempts)
			time.Sleep(time.Duration(db.config.RTimeout) * time.Second)
			return db.connect()
		}
		return fmt.Errorf("Failed connection: %v", err)
	}
	db.conn = conn
	return nil
}

func (db *MySQL) checkConnection() error {
	if db.conn == nil {
		return db.connect()
	}
	return nil
}
