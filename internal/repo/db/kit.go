package db

import (
	"database/sql"
	"fmt"
	"strings"

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
		// TODO: move to mapper
		if kit.Tags.Valid {
			for v := range strings.SplitSeq(kit.Tags.String, ",") {
				kit.tagList = append(kit.tagList, KitTag{Kit: kit.Id, Name: v})
			}
		}
		kits = append(kits, *DbToKit(&kit))
	}
	return &kits, nil
}

// TODO: ON CONFLICT UPDATE
func (d *Sqlite) StoreKit(kit *m.Kit) (kitId int64, err error) {
	kitdb := kitToDb(kit)
	sql := `insert into kit(uid, name, iscustom, description, copyright, licence, credits, url) values(:uid, :name, :iscustom, :description, :copyright, :licence, :credits, :url)`

	tx, err := d.Db.Beginx()
	if err != nil {
		return 0, fmt.Errorf("failed store kit: %w", err)
	}

	// insert kit
	res, err := tx.NamedExec(sql, kitdb)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed store kit: %w", err)
	}
	kitId, err = res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return kitId, fmt.Errorf("failed store kit: %w", err)
	}

	// insert tags
	if len(kitdb.tagList) != 0 {
		for i := range kitdb.tagList {
			kitdb.tagList[i].Kit = kitId
		}

		res, err = tx.NamedExec("insert into kit_tag(kit, name) values(:kit, :name)", kitdb.tagList)
		if err != nil {
			tx.Rollback()
			return kitId, fmt.Errorf("failed store kit: %w", err)
		}
	}
	tx.Commit()

	return kitId, nil
}
