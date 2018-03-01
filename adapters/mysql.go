package adapters

import (
	"database/sql"
	"fmt"

	parser "github.com/PGV65/sql-parser"
	_ "github.com/go-sql-driver/mysql" // mysql
)

// MySQL - DB adapter
type MySQL struct {
	config Config
	conn   *sql.DB
	SQLAdapter
	migrations migrations
	processed  map[string]Migration
}

// NewMySQL - MySQL adapter constructor
func NewMySQL(config Config) *MySQL {
	return &MySQL{
		config: config,
	}
}

// AddMigrations - adding migration
func (db *MySQL) AddMigrations(m migrations) {
	db.migrations = m
}

// Check - checking for migrations
func (db *MySQL) Check() (migrations map[string]Migration, err error) {
	migrations = make(map[string]Migration)
	SQL := `SELECT name, migrated_at FROM migrations`
	rows, err := db.exec(SQL)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var m Migration
		if err = rows.Scan(&m.Name, &m.MigratedAt); err != nil {
			return
		}
		migrations[m.Name] = m
		fmt.Println("Rows: ", m)
	}
	db.processed = migrations
	return
}

// Migrate - starting migration
func (db *MySQL) Migrate() (err error) {
	var keys []string
	for key := range db.migrations {
		if _, ok := db.processed[key]; ok == false {
			keys = append(keys, key)
		}
	}
	fmt.Println("Keys: ", keys)
	for _, key := range keys {
		if err = db.migrate(db.migrations[key]); err != nil {
			return
		}
		if err = db.save(key); err != nil {
			return
		}
	}
	fmt.Println("Processed!")
	return
}

func (db *MySQL) save(key string) (err error) {
	SQL := `INSERT INTO migrations (name) VALUES ('` + key + `');`
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
	connectionInfo := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)
	conn, err := sql.Open("mysql", connectionInfo)
	if err != nil {
		fmt.Println("MySQL connection error", err)
		return err
	}
	if conn == nil {
		fmt.Println("Connection to MySQL is nil")
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
