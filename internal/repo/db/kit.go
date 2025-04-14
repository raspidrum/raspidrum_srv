package db

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	m "github.com/raspidrum-srv/internal/model"
)

type KitDb struct {
	Id          int64          `db:"id"`
	Uid         string         `db:"uid"`
	Name        string         `db:"name"`
	IsCustom    int            `db:"iscustom"`
	Description sql.NullString `db:"description"`
	Copyright   sql.NullString `db:"copyright"`
	Licence     sql.NullString `db:"licence"`
	Credits     sql.NullString `db:"credits"`
	Url         sql.NullString `db:"url"`
	Tags        sql.NullString `db:"tags"`
	tagList     []KitTag
}

type KitTag struct {
	Id   int64  `db:"id"`
	Kit  int64  `db:"kit"`
	Name string `db:"name"`
}

// TODO: optional filter by
//   - like name
//   - isCustom
//   - in (tags)
func (d *Sqlite) ListKits() (*[]m.Kit, error) {
	rows, err := d.Db.Queryx(`select k.*, string_agg(t.name, ',') as tags
	from kit k left join kit_tag t on t.kit = k.id
	group by k.id, k.uid, k.name, k.iscustom, k.description, k.copyright, k.licence, k.credits, k.url
	order by k.name, k.id
	`)
	if err != nil {
		return nil, fmt.Errorf("failed sql: %w", err)
	}
	defer rows.Close()
	kits := []m.Kit{}
	for rows.Next() {
		kit := KitDb{}
		err := rows.StructScan(&kit)
		if err != nil {
			return nil, fmt.Errorf("failed sql: %w", err)
		}
		kits = append(kits, *dbToKit(&kit))
	}
	return &kits, nil
}

// TODO: ON CONFLICT UPDATE
func (d *Sqlite) StoreKit(tx *sqlx.Tx, kit *m.Kit) (kitId int64, err error) {
	localTx := tx == nil
	kitdb := kitToDb(kit)
	sql := `insert into kit(uid, name, iscustom, description, copyright, licence, credits, url) values(:uid, :name, :iscustom, :description, :copyright, :licence, :credits, :url)`

	if localTx {
		tx, err = d.Db.Beginx()
		if err != nil {
			return 0, fmt.Errorf("failed store kit: %w", err)
		}
	}

	// insert kit
	res, err := tx.NamedExec(sql, kitdb)
	if err != nil {
		if localTx {
			tx.Rollback()
		}
		return 0, fmt.Errorf("failed store kit: %w", err)
	}
	kitId, err = res.LastInsertId()
	if err != nil {
		if localTx {
			tx.Rollback()
		}
		return kitId, fmt.Errorf("failed store kit: %w", err)
	}

	// insert tags
	if len(kitdb.tagList) != 0 {
		for i := range kitdb.tagList {
			kitdb.tagList[i].Kit = kitId
		}

		res, err = tx.NamedExec("insert into kit_tag(kit, name) values(:kit, :name)", kitdb.tagList)
		if err != nil {
			if localTx {
				tx.Rollback()
			}
			return kitId, fmt.Errorf("failed store kit: %w", err)
		}
	}
	if localTx {
		tx.Commit()
	}

	return kitId, nil
}

func (d *Sqlite) getKitByUid(tx *sqlx.Tx, uids []string, fields ...string) (*map[string]KitDb, error) {
	fs := mapKitFields(fields)
	// add required fields
	fs["uid"] = void{}
	fs["id"] = void{}

	var err error
	localTx := tx == nil
	if localTx {
		tx, err = d.Db.Beginx()
		if err != nil {
			return nil, fmt.Errorf("failed getKitByUid: %w", err)
		}
	}

	sql := fmt.Sprintf("select %s from kit where uid in (?)", flatFieldMap(fs))
	sql, args, err := sqlx.In(sql, uids)
	if err != nil {
		if localTx {
			tx.Rollback()
		}
		return nil, fmt.Errorf("failed getKitByUid: %w", err)
	}
	rows, err := tx.Queryx(sql, args...)
	if err != nil {
		if localTx {
			tx.Rollback()
		}
		return nil, fmt.Errorf("failed getKitByUid: %w", err)
	}
	var dberr error
	defer func() {
		rows.Close()
		if localTx {
			if dberr != nil {
				tx.Rollback()
			} else {
				tx.Commit()
			}
		}
	}()

	res := make(map[string]KitDb, 0)
	for rows.Next() {
		kit := KitDb{}
		dberr = rows.StructScan(&kit)
		if dberr != nil {
			return nil, fmt.Errorf("failed getKitByUid: %w", err)
		}
		res[kit.Uid] = kit
	}

	return &res, nil
}

func mapKitFields(fields []string) fieldMap {
	dbFields := map[string]string{
		"id":          "id",
		"uuid":        "uid",
		"name":        "name",
		"iscustom":    "iscustom",
		"description": "description",
		"copyright":   "copyright",
		"licence":     "licence",
		"credits":     "credits",
		"url":         "url",
	}
	res := make(fieldMap, 0)
	for _, v := range fields {
		if vr, ok := dbFields[v]; ok {
			res[vr] = void{}
		}
	}
	return res
}
