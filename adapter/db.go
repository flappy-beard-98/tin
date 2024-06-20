package adapter

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path/filepath"
)

type Db struct {
	*sqlx.DB
}

func NewSqliteDb(filename string) (*Db, error) {
	err := createDirectoryIfNotExists(filename)
	if err != nil {
		return nil, err
	}
	db, err := sqlx.Open("sqlite3", filename)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &Db{DB: db}, nil
}

func createDirectoryIfNotExists(path string) error {
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func (o *Db) Close() {
	_ = o.DB.Close()
}
