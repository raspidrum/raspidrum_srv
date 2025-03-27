package loadkit

import (
	"fmt"

	u "github.com/google/uuid"

	m "github.com/raspidrum-srv/internal/model"
	db "github.com/raspidrum-srv/internal/repo/db"
	f "github.com/raspidrum-srv/internal/repo/file"
)

func LoadKit(path string, db *db.Sqlite) (kitId int64, err error) {
	items, err := f.ParseYAMLDir(path)
	if err != nil {
		return kitId, fmt.Errorf("failed load kit files: %w", err)
	}

	// search kit
	kitCnt := 0
	var kitKey string
	for k, v := range items {
		switch v.(type) {
		case *m.Instrument:
			continue
		case *m.Kit:
			kitCnt++
			kitKey = k
		default:
			return kitId, fmt.Errorf("unknown format in: %s %w", k, err)
		}
	}
	// check for one kit in path
	if kitCnt > 1 {
		return kitId, fmt.Errorf("too may kit-files: %d. must be one", kitCnt)
	}
	if len(kitKey) == 0 {
		return kitId, fmt.Errorf("not found kit file")
	}

	// TODO: store kit and instrument in one transaction

	// store kit
	uuid, err := u.NewV7()
	if err != nil {
		return kitId, fmt.Errorf("failed gen uuid for kit: %w", err)
	}
	kit := items[kitKey].(*m.Kit)
	kit.Uid = uuid.String()
	kitId, err = db.StoreKit(kit)
	if err != nil {
		return kitId, err
	}

	// store instruments
	for k, v := range items {
		// skip kit
		if k == kitKey {
			continue
		}
		uuid, err = u.NewV7()
		if err != nil {
			return kitId, fmt.Errorf("failed gen uuid for instrument: %w", err)
		}
		insrt := v.(*m.Instrument)
		insrt.Uid = uuid.String()
		_, err = db.StoreInstrument(kitId, insrt)
		if err != nil {
			return kitId, err
		}
	}
	return kitId, nil
}
