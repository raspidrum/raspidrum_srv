package db

import (
	"database/sql"
	"fmt"
	"strings"

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
func (d *Sqlite) ListInstruments() (*[]m.Instrument, error) {
	rows, err := d.Db.Queryx(`select i.*, string_agg(t.name, ',') as tags
	from instrument i left join instrument_tag t on t.instrument = i.id
	group by i.id, i.uid, i.key, i.name, i.fullname, i.type, i.subtype, i.description, i.copyright, i.licence, i.credits
	order by i.name, i.id`)
	if err != nil {
		return nil, fmt.Errorf("failed sql: %w", err)
	}
	defer rows.Close()
	ins := []m.Instrument{}
	for rows.Next() {
		instr := Instr{}
		err := rows.StructScan(&instr)
		if err != nil {
			return nil, fmt.Errorf("failed sql: %w", err)
		}
		// TODO: move to mapper
		if instr.Tags.Valid {
			for v := range strings.SplitSeq(instr.Tags.String, ",") {
				instr.tagList = append(instr.tagList, InstrTag{Instrument: instr.Id, Name: v})
			}
		}
		ins = append(ins, *DbToInstrument(&instr))
	}
	return &ins, nil
}

// TODO: ON CONFLICT UPDATE
func (d *Sqlite) StoreInstrument(kitId int64, instr *m.Instrument) (instrId int64, err error) {
	instrdb := instrumentToDb(instr)
	sql := `insert into instrument(uid, key, name, fullname, type, subtype, description, copyright, licence, credits, controls, layers)
	values (:uid, :key, :name, :fullname, :type, :subtype, :description, :copyright, :licence, :credits, :controls, :layers)`

	tx, err := d.Db.Beginx()
	if err != nil {
		return 0, fmt.Errorf("failed store instrument: %w", err)
	}

	// insert instrument
	res, err := tx.NamedExec(sql, instrdb)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed store instrument: %w", err)
	}
	instrId, err = res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return instrId, fmt.Errorf("failed store instrument: %w", err)
	}

	// link instrument with kit
	res, err = tx.Exec("insert into kit_instrument(kit, instrument) values(:kit, :instr)", kitId, instrId)
	if err != nil {
		tx.Rollback()
		return instrId, fmt.Errorf("failed store instrument: %w", err)
	}

	// insert tags
	if len(instrdb.tagList) != 0 {
		for i := range instrdb.tagList {
			instrdb.tagList[i].Instrument = instrId
		}

		res, err = tx.NamedExec("insert into instrument_tag(instrument, name) values(:instrument, :name)", instrdb.tagList)
		if err != nil {
			tx.Rollback()
			return instrId, fmt.Errorf("failed store instrument tags: %w", err)
		}
	}
	tx.Commit()

	return instrId, nil
}
