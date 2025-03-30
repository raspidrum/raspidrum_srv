package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"

	m "github.com/raspidrum-srv/internal/model"
)

func kitToDb(kit *m.Kit) *KitDb {
	tgl := make([]KitTag, len(kit.Tags))
	for i, v := range kit.Tags {
		tgl[i] = KitTag{Name: v}
	}
	return &KitDb{
		Id:          kit.Id,
		Uid:         kit.Uid,
		Name:        kit.Name,
		IsCustom:    1,
		Description: sql.NullString{Valid: true, String: kit.Description},
		Copyright:   sql.NullString{Valid: true, String: kit.Copyright},
		Licence:     sql.NullString{Valid: true, String: kit.Licence},
		Credits:     sql.NullString{Valid: true, String: kit.Credits},
		Url:         sql.NullString{Valid: true, String: kit.Url},
		tagList:     tgl,
	}
}

func DbToKit(kit *KitDb) *m.Kit {
	tgl := make([]string, len(kit.tagList))
	for i, v := range kit.tagList {
		tgl[i] = v.Name
	}
	res := m.Kit{
		Id:   kit.Id,
		Uid:  kit.Uid,
		Name: kit.Name,
		Tags: tgl,
	}
	if kit.Description.Valid {
		res.Description = kit.Description.String
	}
	if kit.Copyright.Valid {
		res.Copyright = kit.Copyright.String
	}
	if kit.Licence.Valid {
		res.Licence = kit.Licence.String
	}
	if kit.Credits.Valid {
		res.Credits = kit.Credits.String
	}
	if kit.Url.Valid {
		res.Url = kit.Url.String
	}
	return &res
}

func instrumentToDb(instr *m.Instrument) *Instr {
	// map tags
	tgl := make([]InstrTag, len(instr.Tags))
	for i, v := range instr.Tags {
		tgl[i] = InstrTag{Name: v}
	}

	res := Instr{
		Id:          instr.Id,
		Uid:         instr.Uid,
		Key:         instr.InstrumentKey,
		Name:        instr.Name,
		Fullname:    sql.NullString{Valid: true, String: instr.FullName},
		Type:        instr.Type,
		Subtype:     instr.SubType,
		Description: sql.NullString{Valid: true, String: instr.Description},
		Copyright:   sql.NullString{Valid: true, String: instr.Copyright},
		Licence:     sql.NullString{Valid: true, String: instr.Licence},
		Credits:     sql.NullString{Valid: true, String: instr.Credits},
		tagList:     tgl,
	}

	// convert controls to json
	if len(instr.Controls) > 0 {
		ctrls, err := json.Marshal(instr.Controls)
		if err != nil {
			slog.Error(fmt.Sprint(fmt.Errorf("failed convert to json instrument controls due storing to db: %w", err)))
		}
		res.Controls = sql.NullString{Valid: true, String: string(ctrls)}
	}

	// convert layers to json
	if len(instr.Layers) > 0 {
		lrs, err := json.Marshal(instr.Layers)
		if err != nil {
			slog.Error(fmt.Sprint(fmt.Errorf("failed convert to json instrument layers due storing to db: %w", err)))
		}
		res.Layers = sql.NullString{Valid: true, String: string(lrs)}
	}

	return &res
}

func DbToInstrument(ins *Instr) *m.Instrument {
	tgl := make([]string, len(ins.tagList))
	for i, v := range ins.tagList {
		tgl[i] = v.Name
	}
	res := m.Instrument{
		Id:            ins.Id,
		Uid:           ins.Uid,
		InstrumentKey: ins.Key,
		Name:          ins.Name,
		Type:          ins.Type,
		SubType:       ins.Subtype,
		Tags:          tgl,
	}
	if ins.Fullname.Valid {
		res.FullName = ins.Fullname.String
	}
	if ins.Description.Valid {
		res.Description = ins.Description.String
	}
	if ins.Copyright.Valid {
		res.Copyright = ins.Copyright.String
	}
	if ins.Licence.Valid {
		res.Licence = ins.Licence.String
	}
	if ins.Credits.Valid {
		res.Credits = ins.Credits.String
	}

	if ins.Controls.Valid && len(ins.Controls.String) > 0 {
		var ctrls []m.Controls
		err := json.Unmarshal([]byte(ins.Controls.String), &ctrls)
		if err != nil {
			slog.Error(fmt.Sprint(fmt.Errorf("failed convert instrument controls from json due loading from db: %w", err)))
		}
		res.Controls = ctrls
	}

	if ins.Layers.Valid && len(ins.Layers.String) > 0 {
		var lrs []m.Layer
		err := json.Unmarshal([]byte(ins.Layers.String), &lrs)
		if err != nil {
			slog.Error(fmt.Sprint(fmt.Errorf("failed convert instrument layers from json due loading from db: %w", err)))
		}
		res.Layers = lrs
	}

	return &res

}
