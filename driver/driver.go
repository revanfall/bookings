package driver

import (
	"database/sql"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"time"
)

// DB holds database connection pool
type DB struct {
	Sql *sql.DB
}

var dbConn = &DB{}

const maxOpenDbConn = 10
const maxIdle = 5
const maxDbLifetime = 5 * time.Minute

// ConnectSQL Creates database pool for postgres
func ConnectSQL(dsn string) (*DB, error) {
	d, err := NewDatabase(dsn)
	if err != nil {
		panic(err)
	}
	d.SetMaxOpenConns(maxOpenDbConn)
	d.SetConnMaxIdleTime(maxIdle)
	d.SetConnMaxLifetime(maxDbLifetime)

	dbConn.Sql = d

	err = testDB(d)
	if err != nil {
		return nil, err
	}
	return dbConn, nil
}

// tries to ping the database
func testDB(db *sql.DB) error {
	err := db.Ping()
	if err != nil {
		return err
	}
	return nil
}

// NewDatabase creates the new database for application
func NewDatabase(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil

}
