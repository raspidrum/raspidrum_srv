package db

import (
	"database/sql"
	"fmt"
	"strings"
)

// Field Tags MUST NOT be used outside of this package
type Instrument struct {
	Id          int64          `db:"id"`
	Uid         string         `db:"uid"`
	Key         string         `db:"key"`
	Name        string         `db:"name"`
	Fullname    sql.NullString `db:"fullname"`
	Type        string         `db:"type"`
	Subtype     string         `db:"subtype"`
	Description sql.NullString `db:"description"`
	Copyright   sql.NullString `db:"copyright"`
	Licence     sql.NullString `db:"licence"`
	Credits     sql.NullString `db:"credits"`
	Tags        sql.NullString `db:"tags"`
	TagList     []string
}

type InstrumentTag struct {
	Id         int64  `db:"id"`
	Instrument int64  `db:"instrument"`
	Name       string `db:"name"`
}

type kitInstrument struct {
	kit        int64 `db:"kit"`
	instrument int64 `db:"instrument"`
}

// TODO: optional filter by
//   - like name
//   - type, subtype
//   - in (tags)
//   - kit
func (d *Sqlite) ListInstruments() (*[]Instrument, error) {
	ins := []Instrument{}

	rows, err := d.db.Queryx(`select i.*, string_agg(t.name, ',') as tags
	from instrument i left join instrument_tag t on t.instrument = i.id
	group by i.id, i.uid, i.key, i.name, i.fullname, i.type, i.subtype, i.description, i.copyright, i.licence, i.credits
	order by i.name, i.id`)
	if err != nil {
		return &ins, fmt.Errorf("failed sql: %w", err)
	}
	for rows.Next() {
		instr := Instrument{}
		err := rows.StructScan(instr)
		if err != nil {
			return &ins, fmt.Errorf("failed sql: %w", err)
		}
		if instr.Tags.Valid {
			instr.TagList = strings.Split(instr.Tags.String, ",")
			// Field Tags not usable outside from this package
			instr.Tags.Valid = false
		}
		ins = append(ins, instr)
	}
	return &ins, nil
}

// TODO: ON CONFLICT UPDATE
func (d *Sqlite) StoreInstrument(kitId int64, instr Instrument) (instrId int64, err error) {
	sql := `insert into instrument(id, uid, key, name, fullname, type, subtype, description, copyright, licence, credits)
	values (:id, :uid, :key, :name, :fullname, :type, :subtype, :description, :copyright, :licence, :credits)`

	tx, err := d.db.Beginx()
	if err != nil {
		return 0, fmt.Errorf("failed store instrument: %w", err)
	}

	// insert instrument
	res, err := tx.NamedExec(sql, instr)
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
	if len(instr.TagList) != 0 {
		tags := make([]map[string]interface{}, len(instr.TagList))
		for i, v := range instr.TagList {
			tags[i] = map[string]interface{}{"instrument": instrId, "name": v}
		}

		res, err = tx.NamedExec("insert into instrument_tag(instrument, name) values(:instr, :name)", tags)
		if err != nil {
			tx.Rollback()
			return instrId, fmt.Errorf("failed store instrument: %w", err)
		}
	}
	tx.Commit()

	return instrId, nil
}
