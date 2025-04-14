package preset

import (
	"fmt"

	u "github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/raspidrum-srv/internal/repo/db"
	f "github.com/raspidrum-srv/internal/repo/file"
)

func ImportPresetFromFile(path string, db *db.Sqlite) (presetId int64, err error) {
	pst, err := f.ParsePreset(path)
	if err != nil {
		return 0, fmt.Errorf("faild load preset from file: %s %w", path, err)
	}
	// store kit and instrument in one transaction
	err = db.RunInTx(func(tx *sqlx.Tx) error {
		if len(pst.Uid) == 0 {
			uuid, err := u.NewV7()
			if err != nil {
				return fmt.Errorf("failed gen uuid for preset: %w", err)
			}
			pst.Uid = uuid.String()
		}
		presetId, err = db.StorePreset(tx, pst)
		if err != nil {
			return err
		}
		return nil
	})
	return presetId, err
}
