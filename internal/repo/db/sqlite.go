package db

import (
	"fmt"
	"path"

	"github.com/jmoiron/sqlx"
)

const dbname = "kits.sqlite3"

type Sqlite struct {
	db *sqlx.DB
}

func (d *Sqlite) Connect(dbpath string) error {
	var err error

	dbfile := path.Join(dbpath, dbname)
	constr := fmt.Sprintf("file:%s?_foreign_keys=true", dbfile)
	d.db, err = sqlx.Connect("sqlite3", constr)
	if err != nil {
		return fmt.Errorf("failed connect to db: %s %w", constr, err)
	}
	return nil
}
