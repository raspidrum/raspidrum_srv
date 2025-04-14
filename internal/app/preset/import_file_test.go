package preset

import (
	"path"
	"runtime"
	"testing"

	"github.com/raspidrum-srv/internal/repo/db"
)

func getDBPath() string {
	_, f, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(f), "../../../db/")
	return dir
}

func getProjectPath() string {
	_, f, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(f), "../../../")
	return dir
}

func TestImportPresetFromFile(t *testing.T) {
	d := &db.Sqlite{}

	dir := getDBPath()
	err := d.Connect(dir)
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	defer d.Db.Close()

	testDataPath := path.Join(getProjectPath(), "testdata/")

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
