package db

import (
	"reflect"
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

func TestSqlite_getKitByUid(t *testing.T) {
	d := &Sqlite{}
	dir := getDBPath()
	err := d.Connect(dir)
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	defer d.Db.Close()

	type args struct {
		uids   []string
		fields []string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]KitDb
		wantErr bool
	}{
		{
			name: "existing uuid, no specify fields",
			args: args{
				uids: []string{"0195efff-eb8e-78d9-9be3-bf8dde7bbf0a"},
			},
			want: map[string]KitDb{
				"0195efff-eb8e-78d9-9be3-bf8dde7bbf0a": {
					Id:  1,
					Uid: "0195efff-eb8e-78d9-9be3-bf8dde7bbf0a",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := d.getKitByUid(nil, tt.args.uids, tt.args.fields...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Sqlite.getKitByUid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, &tt.want) {
				t.Errorf("Sqlite.getKitByUid() = %v, want %v", got, tt.want)
			}
		})
	}
}
