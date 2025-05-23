package db

import (
	"errors"
	"fmt"
	"path"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type void struct{}
type fieldMap map[string]void

type Sqlite struct {
	db *sqlx.DB
}

func NewSqlite(dbPath string) (*Sqlite, error) {
	var err error

	connStr := fmt.Sprintf("file:%s?_foreign_keys=true", path.Join(dbPath, "kits.sqlite3"))
	db, err := sqlx.Connect("sqlite3", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed connect to db: %s %w", connStr, err)
	}
	return &Sqlite{db: db}, nil
}

func (d *Sqlite) RunInTx(fn func(tx *sqlx.Tx) error) error {
	tx, err := d.db.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			err = errors.Join(err, rollbackErr)
		}
	}()

	err = fn(tx)
	if err == nil {
		return tx.Commit()
	}

	return err
}

// Close closes the database connection
func (s *Sqlite) Close() error {
	return s.db.Close()
}

func flatFieldMap(fs fieldMap) string {
	fss := make([]string, len(fs))
	i := 0
	for k := range fs {
		fss[i] = k
		i++
	}
	return strings.Join(fss, ",")
}
