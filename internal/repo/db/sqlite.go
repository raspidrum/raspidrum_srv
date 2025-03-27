package db

import (
	"fmt"
	"path"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const dbname = "kits.sqlite3"

type Sqlite struct {
	Db *sqlx.DB
}

func (d *Sqlite) Connect(dbpath string) error {
	var err error

	dbfile := path.Join(dbpath, dbname)
	constr := fmt.Sprintf("file:%s?_foreign_keys=true", dbfile)
	d.Db, err = sqlx.Connect("sqlite3", constr)
	if err != nil {
		return fmt.Errorf("failed connect to db: %s %w", constr, err)
	}
	return nil
}
