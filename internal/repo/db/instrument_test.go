package db

import (
	"reflect"
	"testing"
)

func TestSqlite_ListInstruments(t *testing.T) {
	dir := getDBPath()
	d, err := NewSqlite(dir)
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	defer d.Close()

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

func Test_mapInstrFields(t *testing.T) {

	tests := []struct {
		name   string
		fields []string
		want   fieldMap
	}{
		{
			name:   "all known fields",
			fields: []string{"id", "uuid", "name"},
			want:   map[string]void{"id": {}, "uid": {}, "name": {}},
		},
		{
			name:   "ukknown fields",
			fields: []string{"id", "uuid", "name", "id from dual; delete from kit;"},
			want:   map[string]void{"id": {}, "uid": {}, "name": {}},
		},
		{
			name:   "empty list",
			fields: []string{},
			want:   map[string]void{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mapInstrFields(tt.fields); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mapInstrFields() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSqlite_getInstrumentsByUid(t *testing.T) {
	d, err := NewSqlite(getDBPath())
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	defer d.Close()

	type args struct {
		uids   []string
		fields []string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]Instr
		wantErr bool
	}{
		{
			name: "some exisiting uuids, no specify fields",
			args: args{
				uids:   []string{"0195efff-eb8f-72af-847e-e9322cefe203", "0195efff-eb8f-77e3-8da6-dec7da389977"},
				fields: []string{},
			},
			want: map[string]Instr{
				"0195efff-eb8f-72af-847e-e9322cefe203": {
					Id:  2,
					Uid: "0195efff-eb8f-72af-847e-e9322cefe203",
				},
				"0195efff-eb8f-77e3-8da6-dec7da389977": {
					Id:  17,
					Uid: "0195efff-eb8f-77e3-8da6-dec7da389977",
				},
			},
		},
		{
			name: "some exisiting uuids, specify required fields",
			args: args{
				uids:   []string{"0195efff-eb8f-72af-847e-e9322cefe203", "0195efff-eb8f-77e3-8da6-dec7da389977"},
				fields: []string{"id", "uuid", "key"},
			},
			want: map[string]Instr{
				"0195efff-eb8f-72af-847e-e9322cefe203": {
					Id:  2,
					Uid: "0195efff-eb8f-72af-847e-e9322cefe203",
					Key: "hihat",
				},
				"0195efff-eb8f-77e3-8da6-dec7da389977": {
					Id:  17,
					Uid: "0195efff-eb8f-77e3-8da6-dec7da389977",
					Key: "tom2",
				},
			},
		},

		{
			name: "some not exisiting uuids, specify required fields",
			args: args{
				uids:   []string{"0195efff-eb8f-72af-847e-e9322cefe203", "0195efff-eb8f-77e3-8da6-dec7da389977", "aaaa-bbbb"},
				fields: []string{"id", "uuid"},
			},
			want: map[string]Instr{
				"0195efff-eb8f-72af-847e-e9322cefe203": {
					Id:  2,
					Uid: "0195efff-eb8f-72af-847e-e9322cefe203",
				},
				"0195efff-eb8f-77e3-8da6-dec7da389977": {
					Id:  17,
					Uid: "0195efff-eb8f-77e3-8da6-dec7da389977",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := d.getInstrumentsByUid(nil, tt.args.uids, tt.args.fields...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Sqlite.getInstrumentsByUid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, &tt.want) {
				t.Errorf("Sqlite.getInstrumentsByUid() = %v, want %v", got, tt.want)
			}
		})
	}
}
