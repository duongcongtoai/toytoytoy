package common

import (
	"context"
	"database/sql"
	"time"

	"github.com/duongcongtoai/toytoytoy/sqlc/togo"
)

type DBX interface {
	togo.DBTX
	BeginTx(ctx context.Context, opts *sql.TxOptions) (Tx, error)
	Begin() (Tx, error)
}

type Tx interface {
	Commit() error
	Rollback() error
	togo.DBTX
}
type SqlDB struct {
	*sql.DB
}

func (db *SqlDB) BeginTx(ctx context.Context, opts *sql.TxOptions) (Tx, error) {
	return db.DB.BeginTx(ctx, opts)
}
func (db *SqlDB) Begin() (Tx, error) {
	return db.DB.Begin()
}

type Config struct {
	DSN             string
	MaxConnIdleTime time.Duration
	MaxIdleConn     int
	MaxOpenConn     int
}

func ConnectDB(c Config) *SqlDB {
	db, err := sql.Open("mysql", c.DSN)
	if err != nil {
		panic(err)
	}
	db.SetConnMaxIdleTime(c.MaxConnIdleTime)
	db.SetMaxIdleConns(c.MaxIdleConn)
	db.SetMaxOpenConns(c.MaxOpenConn)
	return &SqlDB{db}
}
func CleanUpTestData(db *SqlDB) error {
	_, err := db.Exec("TRUNCATE TABLE wagers")
	if err != nil {
		return err
	}
	_, err = db.Exec("TRUNCATE TABLE purchases")
	if err != nil {
		return err
	}
	return nil
}
