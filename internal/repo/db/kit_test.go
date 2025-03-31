package db

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestSqlite_ListKits(t *testing.T) {

	d := &Sqlite{}

	dir := getDBPath()
	err := d.Connect(dir)
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	defer d.Db.Close()

	tests := []struct {
		name string
		//want    *[]KitDb
		wantLen int
		wantErr bool
	}{
		{
			name: "all kits",
			//want:    nil,
			wantLen: 1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := d.ListKits()
			if (err != nil) != tt.wantErr {
				t.Errorf("Sqlite.ListKits() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantLen != -1 && len(*got) != tt.wantLen {
				t.Errorf("Sqlite.ListKits() len = %v, want len = %v", len(*got), tt.wantLen)
			}
			//if tt.wantLen != 0 && !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("Sqlite.ListKits() = %v, want %v", got, tt.want)
			//}
		})
	}
}
