package db

import (
	"reflect"
	"testing"
)

func TestSqlite_ListInstruments(t *testing.T) {
	d := &Sqlite{}

	dir := getRootPath()
	err := d.Connect(dir)
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	defer d.db.Close()

	tests := []struct {
		name    string
		want    *[]Instrument
		wantLen int
		wantErr bool
	}{
		{
			name:    "empty list",
			want:    nil,
			wantLen: 0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := d.ListInstruments()
			if (err != nil) != tt.wantErr {
				t.Errorf("Sqlite.ListInstruments() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantLen != -1 && len(*got) != tt.wantLen {
				t.Errorf("Sqlite.ListInstruments() len = %v, want len = %v", len(*got), tt.wantLen)
			}
			if tt.wantLen != 0 && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Sqlite.ListInstruments() = %v, want %v", got, tt.want)
			}
		})
	}
}
