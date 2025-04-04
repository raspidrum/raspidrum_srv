package db

import "database/sql"

type KitPrst struct {
	Id          int64  `db:"id"`
	Uid         string `db:"uid"`
	KitUid      string
	KitId       int64  `db:"kit"`
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

type PrtsInstr struct {
	Id         int64  `db:"id"`
	PresetId   int64  `db:"preset"`
	ChannelId  int64  `db:"channel"`
	ChannelKey string `db:"channel_key"`
	InstrId    int64  `db:"instrument"`
	InstrUid   string
	Name       string         `db:"name"`
	MidiKey    sql.NullString `db:"midikey"`
	Controls   string         `db:"controls"`
	Layers     sql.NullString `db:"layers"`
}
