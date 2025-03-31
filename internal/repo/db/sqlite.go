package db

import (
	"fmt"
	"path"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// TODO: move to config
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

// Where conditions
type Condition func() (sql string, args []interface{}, err error)

func Eq(field string, inargs ...interface{}) Condition {
	return func() (sql string, args []interface{}, err error) {
		return fmt.Sprintf("%s = ?", field), inargs, nil
	}
}

func buildConditions(conds ...Condition) (sql string, args []interface{}, err error) {
	sqls := make([]string, len(conds))
	for i, cond := range conds {
		s, a, err := cond()
		if err != nil {
			return "", nil, fmt.Errorf("failed construct sql condition: %w", err)
		}
		sqls[i] = s
		args = append(args, a...)
	}
	if len(sqls) > 0 {
		sql = "where " + strings.Join(sqls, " and ")
	}
	return
}
