package db

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	m "github.com/raspidrum-srv/internal/model"
)

type KitBase struct {
	KitId       int64  `db:"kit"`
	KitUid      string `db:"kit_uid"`
	KitName     string `db:"kit_name"`
	KitIsCustom int    `db:"kit_iscustom"`
}
type KitPrst struct {
	KitBase
	Id          int64  `db:"id"`
	Uid         string `db:"uid"`
	Name        string `db:"name"`
	Channels    []PrstChnl
	Instruments []PrtsInstr
}

type PrstChnl struct {
	Id       int64  `db:"id"`
	PresetId int64  `db:"preset"`
	Key      string `db:"key"`
	Name     string `db:"name"`
	Controls string `db:"controls"`
}

type InstrBase struct {
	InstrId       int64          `db:"instrument"`
	InstrUid      string         `db:"instrument_uid"`
	InstrKey      string         `db:"instrument_key"`
	InstrName     string         `db:"instrument_name"`
	InstrMidiKey  sql.NullString `db:"instrument_midikey"`
	InstrControls sql.NullString `db:"instrument_controls"`
	InstrLayers   sql.NullString `db:"instrument_layers"`
}

type PrtsInstr struct {
	InstrBase
	Id         int64          `db:"id"`
	PresetId   int64          `db:"preset"`
	ChannelId  int64          `db:"channel"`
	ChannelKey string         `db:"channel_key"`
	Name       string         `db:"name"`
	MidiKey    sql.NullString `db:"midikey"`
	Controls   string         `db:"controls"`
	Layers     sql.NullString `db:"layers"`
}

func (d *Sqlite) StorePreset(tx *sqlx.Tx, preset *m.KitPreset) (presetId int64, err error) {
	localTx := tx == nil
	pstDb := kitPresetToDb(preset)

	if localTx {
		tx, err = d.Db.Beginx()
		if err != nil {
			return 0, fmt.Errorf("failed store kit preset: %w", err)
		}
	}

	defer func() {
		if localTx {
			if err != nil {
				tx.Rollback()
			} else {
				tx.Commit()
			}
		}
	}()

	// get kitId by kit uuid
	kitIds, err := d.getKitByUid(tx, []string{pstDb.KitUid})
	if err != nil {
		return presetId, err
	}
	if kitIds == nil || len(*kitIds) == 0 {
		return presetId, fmt.Errorf("not found kit with uuid: %s", pstDb.KitUid)
	}
	pstDb.KitId = (*kitIds)[pstDb.KitUid].Id

	// get instruments id
	insUids := make([]string, len(pstDb.Instruments))
	for i, v := range pstDb.Instruments {
		insUids[i] = v.InstrUid
	}
	insIds, err := d.getInstrumentsByUid(tx, insUids)
	if err != nil {
		return presetId, err
	}
	missingInstrs := []string{}
	for i, v := range pstDb.Instruments {
		instr, ok := (*insIds)[v.InstrUid]
		if !ok {
			missingInstrs = append(missingInstrs, fmt.Sprintf("not found instrument with uuid: %s", v.InstrUid))
			continue
		}
		pstDb.Instruments[i].InstrId = instr.Id
	}
	if len(missingInstrs) != 0 {
		return presetId, fmt.Errorf("failed store kit preset: %s", strings.Join(missingInstrs, "\n"))
	}

	// store kit preset
	sql := `insert into kit_preset(uid, kit, name) values(:uid, :kit, :name)
	on conflict (id) do update set name = excluded.name, uid = excluded.uid
	on conflict (uid) do update set name = excluded.name
	returning id`
	rows, err := tx.NamedQuery(sql, pstDb)
	if err != nil {
		return presetId, fmt.Errorf("failed store kit preset: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&presetId)
		if err != nil {
			return presetId, fmt.Errorf("failed store kit preset: %w", err)
		}
	}

	// store preset channels
	sql = `insert into preset_channel(preset, key, name, controls) values(:preset, :key, :name, :controls) 
	on conflict (id) do update set preset = excluded.preset, key = excluded.key, name = excluded.name, controls = excluded.controls
	on conflict (preset, key) do update set name = excluded.name, controls = excluded.controls
	returning id`
	chnls := make(map[string]int64, len(pstDb.Channels))
	for i, v := range pstDb.Channels {
		pstDb.Channels[i].PresetId = presetId
		rows, err := tx.NamedQuery(sql, pstDb.Channels[i])
		if err != nil {
			return presetId, fmt.Errorf("failed store channel with key: %s of kit preset: %w", v.Key, err)
		}
		defer rows.Close()
		for rows.Next() {
			var chnId int64
			err := rows.Scan(&chnId)
			if err != nil {
				return presetId, fmt.Errorf("failed store channel with key: %s of kit preset: %w", v.Key, err)
			}
			chnls[v.Key] = chnId
		}
	}

	// store preset instruments
	for i, v := range pstDb.Instruments {
		pstDb.Instruments[i].PresetId = presetId
		chnlId, ok := chnls[v.ChannelKey]
		if !ok {
			return presetId, fmt.Errorf("not found channel key: %s for instrument:%s : %w", v.ChannelKey, v.Name, err)
		}
		pstDb.Instruments[i].ChannelId = chnlId
	}
	sql = `insert into preset_instrument(preset, channel, instrument, name, midikey, controls, layers) 
	values(:preset, :channel, :instrument, :name, :midikey, :controls, :layers)
	on conflict (preset, name) do update set channel = excluded.channel, instrument = excluded.instrument, midikey = excluded.midikey, controls = excluded.controls, layers = excluded.layers`
	_, err = tx.NamedExec(sql, pstDb.Instruments)
	if err != nil {
		return presetId, fmt.Errorf("failed store instruments of kit preset: %w", err)
	}

	if localTx {
		tx.Commit()
	}
	return presetId, err
}

// Return minimal list of Kit Presets with minimal info
func (d *Sqlite) ListPresets(conds ...Condition) (*[]m.KitPreset, error) {
	sql_select := `select * from v_kit_preset`
	sql_order := `order by name`
	sql_where, args, err := buildConditions(conds...)
	if err != nil {
		return nil, fmt.Errorf("failed ListPresets: %w", err)
	}
	sql := fmt.Sprintf("%s %s %s", sql_select, sql_where, sql_order)
	rows, err := d.Db.Queryx(sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed ListPresets: %w", err)
	}
	defer rows.Close()
	psts := []m.KitPreset{}
	for rows.Next() {
		pst := KitPrst{}
		err := rows.StructScan(&pst)
		if err != nil {
			return nil, fmt.Errorf("failed ListPresets: %w", err)
		}
		psts = append(psts, *dbToKitPreset(&pst))
	}
	return &psts, nil
}

func (d *Sqlite) GetPreset(conds ...Condition) (*m.KitPreset, error) {
	if len(conds) == 0 {
		return nil, fmt.Errorf("empty conditions")
	}
	// Get Preset overall info
	sql_select := `select * from v_kit_preset`
	sql_where, args, err := buildConditions(conds...)
	if err != nil {
		return nil, fmt.Errorf("failed GetPreset: %w", err)
	}
	sqlstmn := fmt.Sprintf("%s %s", sql_select, sql_where)

	pst := KitPrst{}
	err = d.Db.Get(&pst, sqlstmn, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, fmt.Errorf("failed GetPreset: %w", err)
		}
	}

	// Get Preset channels
	sqlstmn = `select * from preset_channel where preset = :preset`
	pst.Channels = []PrstChnl{}
	err = d.Db.Select(&pst.Channels, sqlstmn, &pst.Id)
	if err != nil {
		return nil, fmt.Errorf("failed GetPreset: %w", err)
	}

	// Get Preset instruments
	sqlstmn = `select * from v_preset_instrument where preset = :preset`
	pst.Instruments = []PrtsInstr{}
	err = d.Db.Select(&pst.Instruments, sqlstmn, &pst.Id)
	if err != nil {
		return nil, fmt.Errorf("failed GetPreset: %w", err)
	}

	return dbToKitPreset(&pst), nil
}
