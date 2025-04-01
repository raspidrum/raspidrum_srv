package db

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	m "github.com/raspidrum-srv/internal/model"
)

// Field Tags MUST NOT be used outside of this package
type Instr struct {
	Id          int64          `db:"id"`
	Uid         string         `db:"uid"`
	Key         string         `db:"key"`
	Name        string         `db:"name"`
	Fullname    sql.NullString `db:"fullname"`
	Type        string         `db:"type"`
	Subtype     string         `db:"subtype"`
	MidiKey     sql.NullString `db:"midikey"`
	Description sql.NullString `db:"description"`
	Copyright   sql.NullString `db:"copyright"`
	Licence     sql.NullString `db:"licence"`
	Credits     sql.NullString `db:"credits"`
	Tags        sql.NullString `db:"tags,omitempty"`
	tagList     []InstrTag
	Controls    sql.NullString `db:"controls"`
	Layers      sql.NullString `db:"layers"`
}

type InstrTag struct {
	Id         int64  `db:"id"`
	Instrument int64  `db:"instrument"`
	Name       string `db:"name"`
}

// TODO: optional filter by
//   - like name
//   - type, subtype
//   - in (tags)
//   - kit
func (d *Sqlite) ListInstruments(conds ...Condition) (*[]m.Instrument, error) {
	sql_select := `select i.*, string_agg(t.name, ',') as tags
	from instrument i join kit_instrument ki on ki.instrument = i.id
	     left join instrument_tag t on t.instrument = i.id`

	sql_group := "group by i.id, i.uid, i.key, i.name, i.fullname, i.type, i.subtype, i.description, i.copyright, i.licence, i.credits"
	sql_order := "order by i.name, i.id"

	sql_where, args, err := buildConditions(conds...)
	if err != nil {
		return nil, fmt.Errorf("failed ListInstruments: %w", err)
	}

	sql := fmt.Sprintf("%s %s %s %s", sql_select, sql_where, sql_group, sql_order)

	rows, err := d.Db.Queryx(sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed sqListInstrumentsl: %w", err)
	}
	defer rows.Close()
	ins := []m.Instrument{}
	for rows.Next() {
		instr := Instr{}
		err := rows.StructScan(&instr)
		if err != nil {
			return nil, fmt.Errorf("failed ListInstruments: %w", err)
		}
		ins = append(ins, *DbToInstrument(&instr))
	}
	return &ins, nil
}

func ByKitId(kit int) Condition {
	return Eq("kit", kit)
}

// TODO: ON CONFLICT UPDATE
func (d *Sqlite) StoreInstrument(tx *sqlx.Tx, kitId int64, instr *m.Instrument) (instrId int64, err error) {
	localTx := tx == nil
	instrdb := instrumentToDb(instr)
	sql := `insert into instrument(uid, key, name, fullname, type, subtype, description, copyright, licence, credits, controls, layers)
	values (:uid, :key, :name, :fullname, :type, :subtype, :description, :copyright, :licence, :credits, :controls, :layers)`

	if localTx {
		tx, err = d.Db.Beginx()
		if err != nil {
			return 0, fmt.Errorf("failed store kit: %w", err)
		}
	}

	// insert instrument
	res, err := tx.NamedExec(sql, instrdb)
	if err != nil {
		if localTx {
			tx.Rollback()
		}
		return 0, fmt.Errorf("failed store instrument: %w", err)
	}
	instrId, err = res.LastInsertId()
	if err != nil {
		if localTx {
			tx.Rollback()
		}
		return instrId, fmt.Errorf("failed store instrument: %w", err)
	}

	// link instrument with kit
	res, err = tx.Exec("insert into kit_instrument(kit, instrument) values(:kit, :instr)", kitId, instrId)
	if err != nil {
		if localTx {
			tx.Rollback()
		}
		return instrId, fmt.Errorf("failed store instrument: %w", err)
	}

	// insert tags
	if len(instrdb.tagList) != 0 {
		for i := range instrdb.tagList {
			instrdb.tagList[i].Instrument = instrId
		}

		res, err = tx.NamedExec("insert into instrument_tag(instrument, name) values(:instrument, :name)", instrdb.tagList)
		if err != nil {
			if localTx {
				tx.Rollback()
			}
			return instrId, fmt.Errorf("failed store instrument tags: %w", err)
		}
	}
	if localTx {
		tx.Commit()
	}

	return instrId, nil
}
