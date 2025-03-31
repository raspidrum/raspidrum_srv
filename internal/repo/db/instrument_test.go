package db

import (
	"testing"
)

func TestSqlite_ListInstruments(t *testing.T) {
	d := &Sqlite{}

	dir := getDBPath()
	err := d.Connect(dir)
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	defer d.Db.Close()

	type args struct {
		conds []Condition
	}
	tests := []struct {
		name string
		args args
		//want    *[]m.Instrument
		wantLen int
		wantErr bool
	}{
		{
			name: "list for kit",
			args: args{
				conds: []Condition{
					ByKitId(1),
				},
			},
			//want:    nil,
			wantLen: 22,
			wantErr: false,
		},
		{
			name: "list for non exists kit",
			args: args{
				conds: []Condition{
					ByKitId(-1),
				},
			},
			//want:    nil,
			wantLen: 0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := d.ListInstruments(tt.args.conds...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Sqlite.ListInstruments() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantLen != -1 && len(*got) != tt.wantLen {
				t.Errorf("Sqlite.ListInstruments() len = %v, want len = %v", len(*got), tt.wantLen)
			}
			//if tt.wantLen != 0 && !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("Sqlite.ListInstruments() = %v, want %v", got, tt.want)
			//}
		})
	}
}
