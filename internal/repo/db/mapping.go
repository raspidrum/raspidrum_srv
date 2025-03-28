package db

import (
	"database/sql"

	m "github.com/raspidrum-srv/internal/model"
)

func kitToDb(kit *m.Kit) *KitDb {
	tgl := make([]string, len(kit.Tags))
	for i := 0; i < len(kit.Tags); i++ {
		tgl[i] = kit.Tags[i]
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
		TagList:     tgl,
	}
}

func instrumentToDb(instr *m.Instrument) *Instr {
	tgl := make([]string, len(instr.Tags))
	for i := 0; i < len(instr.Tags); i++ {
		tgl[i] = instr.Tags[i]
	}
	return &Instr{
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
}
