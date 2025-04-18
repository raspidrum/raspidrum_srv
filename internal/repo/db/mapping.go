package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

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

func dbToKit(kit *KitDb) *m.Kit {
	var tgl []string
	if kit.Tags.Valid {
		tgl = strings.Split(kit.Tags.String, ",")
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
		Key:         instr.Key,
		Name:        instr.Name,
		Fullname:    sql.NullString{Valid: true, String: instr.FullName},
		Type:        instr.Type,
		Subtype:     instr.SubType,
		MidiKey:     sql.NullString{Valid: true, String: instr.MidiKey},
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

func dbToInstrument(ins *Instr) *m.Instrument {
	var tgl []string
	if ins.Tags.Valid {
		tgl = strings.Split(ins.Tags.String, ",")
	}

	res := m.Instrument{
		Id:      ins.Id,
		Uid:     ins.Uid,
		Key:     ins.Key,
		Name:    ins.Name,
		Type:    ins.Type,
		SubType: ins.Subtype,
		Tags:    tgl,
	}
	if ins.Fullname.Valid {
		res.FullName = ins.Fullname.String
	}
	if ins.MidiKey.Valid {
		res.MidiKey = ins.MidiKey.String
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
		var ctrls map[string]m.Control
		err := json.Unmarshal([]byte(ins.Controls.String), &ctrls)
		if err != nil {
			slog.Error(fmt.Sprint(fmt.Errorf("failed convert instrument controls from json due loading from db: %w", err)))
		}
		res.Controls = ctrls
	}

	if ins.Layers.Valid && len(ins.Layers.String) > 0 {
		var lrs map[string]m.Layer
		err := json.Unmarshal([]byte(ins.Layers.String), &lrs)
		if err != nil {
			slog.Error(fmt.Sprint(fmt.Errorf("failed convert instrument layers from json due loading from db: %w", err)))
		}
		res.Layers = lrs
	}

	return &res

}

func kitPresetToDb(pst *m.KitPreset) *KitPrst {
	res := KitPrst{
		Uid: pst.Uid,
		KitBase: KitBase{
			KitUid: pst.Kit.Uid,
		},
		Name: pst.Name,
	}
	// channels
	chs := make([]PrstChnl, len(pst.Channels))
	for i, v := range pst.Channels {
		chs[i] = PrstChnl{
			Key:  v.Key,
			Name: v.Name,
		}
		// marshal controls to json
		if len(v.Controls) > 0 {
			ctrs, err := json.Marshal(v.Controls)
			if err != nil {
				slog.Error(fmt.Sprint(fmt.Errorf("failed convert to json channel controls due storing to db: %w", err)))
			}
			chs[i].Controls = string(ctrs)
		}
	}
	res.Channels = chs

	// instruments
	ins := make([]PrtsInstr, len(pst.Instruments))
	for i, v := range pst.Instruments {
		ins[i] = PrtsInstr{
			InstrBase: InstrBase{
				InstrUid: v.Instrument.Uid,
			},
			Name:       v.Name,
			ChannelKey: v.ChannelKey,
		}
		if len(v.MidiKey) > 0 {
			ins[i].MidiKey = sql.NullString{Valid: true, String: v.MidiKey}
		}
		// marshal controls to json
		if len(v.Controls) > 0 {
			ctrs, err := json.Marshal(v.Controls)
			if err != nil {
				slog.Error(fmt.Sprint(fmt.Errorf("failed convert to json instrument controls due storing to db: %w", err)))
			}
			ins[i].Controls = string(ctrs)
		}
		// marshal layers to json
		if len(v.Layers) > 0 {
			lrs, err := json.Marshal(v.Layers)
			if err != nil {
				slog.Error(fmt.Sprint(fmt.Errorf("failed convert to json instrument layers due storing to db: %w", err)))
			}
			ins[i].Layers = sql.NullString{Valid: true, String: string(lrs)}
		}
	}
	res.Instruments = ins

	return &res
}

func dbToKitPreset(pst *KitPrst) *m.KitPreset {
	res := m.KitPreset{
		Uid: pst.Uid,
		Kit: m.KitRef{
			Id:       pst.KitId,
			Uid:      pst.KitUid,
			Name:     pst.KitName,
			IsCustom: pst.KitIsCustom == 1,
		},
		Name: pst.Name,
	}
	// channels
	chs := make([]m.PresetChannel, len(pst.Channels))
	for i, v := range pst.Channels {
		chs[i] = m.PresetChannel{
			Key:  v.Key,
			Name: v.Name,
		}
		if len(v.Controls) > 0 {
			var ctrls []m.PresetControl
			err := json.Unmarshal([]byte(v.Controls), &ctrls)
			if err != nil {
				slog.Error(fmt.Sprint(fmt.Errorf("failed convert channel controls from json due loading from db: %w", err)))
			}
			chs[i].Controls = ctrls
		}
	}
	res.Channels = chs

	//instruments
	ins := make([]m.PresetInstrument, len(pst.Instruments))
	for i, v := range pst.Instruments {
		ins[i] = m.PresetInstrument{
			Instrument: m.InstrumentRef{
				Id:   v.InstrId,
				Uid:  v.InstrUid,
				Key:  v.InstrKey,
				Name: v.InstrName,
			},
			Name:       v.Name,
			ChannelKey: v.ChannelKey,
		}
		if v.MidiKey.Valid {
			ins[i].MidiKey = v.MidiKey.String
		}
		if len(v.Controls) > 0 {
			var ctrls []m.PresetControl
			err := json.Unmarshal([]byte(v.Controls), &ctrls)
			if err != nil {
				slog.Error(fmt.Sprint(fmt.Errorf("failed convert instrument controls from json due loading from db: %w", err)))
			}
			ins[i].Controls = ctrls
		}
		if v.Layers.Valid && len(v.Layers.String) > 0 {
			var lrs []m.PresetLayer
			err := json.Unmarshal([]byte(v.Layers.String), &lrs)
			if err != nil {
				slog.Error(fmt.Sprint(fmt.Errorf("failed convert instrument layers from json due loading from db: %w", err)))
			}
			ins[i].Layers = lrs
		}
		if v.InstrMidiKey.Valid {
			ins[i].Instrument.MidiKey = v.InstrMidiKey.String
		}
		// Instrument Controls
		if v.InstrControls.Valid && len(v.InstrControls.String) > 0 {
			var ictrls []m.Control
			err := json.Unmarshal([]byte(v.InstrControls.String), &ictrls)
			if err != nil {
				slog.Error(fmt.Sprint(fmt.Errorf("failed convert instrument controls ref from json due loading from db: %w", err)))
			}
			ins[i].Instrument.Controls = ictrls
		}
		// instrument Layers
		if v.InstrLayers.Valid && len(v.InstrLayers.String) > 0 {
			var ilrs []m.Layer
			err := json.Unmarshal([]byte(v.InstrLayers.String), &ilrs)
			if err != nil {
				slog.Error(fmt.Sprint(fmt.Errorf("failed convert instrument layers ref from json due loading from db: %w", err)))
			}
			ins[i].Instrument.Layers = ilrs
		}
	}
	res.Instruments = ins

	return &res
}
