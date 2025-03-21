package db

import (
	"database/sql"
	"fmt"
	"strings"
)

type Instrument struct {
	Id          int            `db:"id"`
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
	Id         int    `db:"id"`
	Instrument int    `db:"instrument"`
	Name       string `db:"name"`
}

type kitInstrument struct {
	kit        int `db:"kit"`
	instrument int `db:"instrument"`
}

// TODO: optional filter by
//   - like name
//   - type, subtype
//   - in (tags)
//   - kit
func (d *Sqlite) ListInstruments() (*[]Instrument, error) {
	ins := []Instrument{}

	rows, err := d.db.Queryx(`select i.*, string_agg(t.name) as tags
	from instrument i left join instrument_tag t on t.instrument = i.id
	group by i.id, i.uid, i.key, i.name, i.fullname, i.type, i.subtype, i.description, i.copyright, i.licence, i.credits
	order by i.name, i.id
	`)
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
		}
		ins = append(ins, instr)
	}
	return &ins, nil
}
