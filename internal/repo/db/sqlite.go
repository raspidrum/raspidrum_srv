package db

import (
	"errors"
	"fmt"
	"path"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// TODO: move to config
const dbname = "kits.sqlite3"

type void struct{}
type fieldMap map[string]void

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

func (d *Sqlite) RunInTx(fn func(tx *sqlx.Tx) error) error {
	tx, err := d.Db.Beginx()
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

func flatFieldMap(fs fieldMap) string {
	fss := make([]string, len(fs))
	i := 0
	for k, _ := range fs {
		fss[i] = k
		i++
	}
	return strings.Join(fss, ",")
}
