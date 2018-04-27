package adapters

import (
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	_ "github.com/lib/pq" // postgres driver
)

const createTablePostgres = "CREATE TABLE %s (name VARCHAR(255) PRIMARY KEY NOT NULL, migrated_at TIMESTAMP NOT NULL); CREATE UNIQUE INDEX _migrations_name_uindex ON _migrations (name)"

// Postgres - DB adapter
type Postgres struct {
	config Config
	SQLAdapter
}

// NewPostgres - Postgres adapter constructor
func NewPostgres(config Config, m migrations) *Postgres {
	config = preparePostgresConfig(config)
	SQLAdapter := SQLAdapter{config: config, migrations: m}
	db := &Postgres{
		config,
		SQLAdapter,
	}
	return db
}

func preparePostgresConfig(config Config) Config {
	if config.RTimeout == 0 {
		config.RTimeout = 1
	}
	if config.RAttempts == 0 {
		config.RAttempts = 1
	}
	if config.User == "" {
		config.User = "postgres"
	}
	if config.Host == "" {
		config.Host = "localhost"
	}
	if config.Port == 0 {
		config.Port = 5432
	}
	if config.SSLmode == "" {
		config.SSLmode = "disable"
	}
	return config
}
func (db *Postgres) getCreateTableQuery() string {
	return fmt.Sprintf(createTablePostgres, table)
}

func (db *Postgres) open(driver, connection string) (*sql.DB, error) {
	return sql.Open(driver, connection)
}

// GetConnectionString
func (db *Postgres) getConnectionString(database string) string {
	return fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=%v",
		db.config.User,
		db.config.Password,
		db.config.Host,
		db.config.Port,
		database,
		db.config.SSLmode,
	)
}

// Migrate - starting migration
func (db *Postgres) Migrate() (err error) {
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
func (db *Postgres) loadProcessed() (err error) {
	var SQL string
	migrations := make(map[string]Migration)

	SQL = `SELECT name, migrated_at FROM ` + table + `;`
	rows, err := db.conn.Query(SQL)
	if err != nil {
		fmt.Println("Error loading migrations from DB: ", err)
		if db.checkError(err) == 2 {
			if err = db.createTable(); err != nil {
				return
			}
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

func (db *Postgres) createTable() (err error) {
	fmt.Println("Creating table: ", table)
	_, err = db.conn.Exec(db.getCreateTableQuery())
	if err != nil {
		fmt.Println("Error creating table: ", err)
		return
	}
	return
}

func (db *Postgres) createDB() (err error) {
	fmt.Println("Creating database: ", db.config.Database)
	db.connect("")
	SQL := `CREATE DATABASE ` + db.config.Database + `;`

	_, err = db.conn.Exec(SQL)
	if err != nil {
		fmt.Println("Error creating DB: ", err)
		return
	}
	db.conn.Close()
	return db.connect(db.config.Database)
}

func (db *Postgres) save(key string) (err error) {
	migratedAt := time.Now().Format("2006-01-02T15:04:05")
	SQL := `INSERT INTO ` + table + ` (name, migrated_at) VALUES ('` + key + `', '` + migratedAt + `');`
	_, err = db.conn.Exec(SQL)
	return
}

func (db *Postgres) migrate(data string) (err error) {
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

func (db *Postgres) connect(database string) (err error) {
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

func (db *Postgres) checkError(err error) (code int) {
	if strings.Contains(err.Error(), fmt.Sprintf("pq: database \"%s\" does not exist", db.config.Database)) {
		return 1
	}
	if strings.Contains(err.Error(), "does not exist") {
		return 2
	}
	return 0
}
