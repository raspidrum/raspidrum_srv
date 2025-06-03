//go:build integration

package preset

import (
	"path"
	"testing"

	db "github.com/raspidrum-srv/internal/repo/db"
)

func TestImportPresetFromFile(t *testing.T) {
	dir := getDBPath()
	d, err := db.NewSqlite(dir)
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	defer d.Close()

	testDataPath := path.Join(getProjectPath(), "testdata/preset_import/")

	tests := []struct {
		name     string
		filename string
		wantErr  bool
	}{
		{
			name:     "preset1",
			filename: "kit_preset_1.yaml",
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ImportPresetFromFile(path.Join(testDataPath, tt.filename), d)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadPresetFromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
