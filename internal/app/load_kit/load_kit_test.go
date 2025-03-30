package loadkit

import (
	"path"
	"testing"

	db "github.com/raspidrum-srv/internal/repo/db"
)

func TestLoadKit(t *testing.T) {
	d := &db.Sqlite{}

	dir := getDBPath()
	err := d.Connect(dir)
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	defer d.Db.Close()

	kitPath := path.Join(getProjectPath(), "../_presets/", "SMDrums")

	tests := []struct {
		name      string
		wantKitId int64
		wantErr   bool
	}{
		{
			name:      "first kit",
			wantKitId: 1,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKitId, err := LoadKit(kitPath, d)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadKit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotKitId != tt.wantKitId {
				t.Errorf("LoadKit() = %v, want %v", gotKitId, tt.wantKitId)
			}
		})
	}
}
