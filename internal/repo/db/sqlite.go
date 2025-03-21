package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

const dbname = "db.sqlite"

type Sqlite struct {
	db *sqlx.DB
}

func (d *Sqlite) Connect() error {
	var err error
	d.db, err = sqlx.Connect("sqlite3", fmt.Sprintf("file:%s?_foreign_keys=true", dbname))
	if err != nil {
		return fmt.Errorf("failed connect to db: %w", err)
	}
	return nil
}
