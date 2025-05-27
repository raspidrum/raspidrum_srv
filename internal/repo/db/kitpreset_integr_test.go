//go:build integration

package db

import (
	"testing"
)

func TestSqlite_GetPreset(t *testing.T) {
	d, err := NewSqlite(getDBPath())
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	defer d.Close()

	type args struct {
		conds []Condition
	}
	tests := []struct {
		name    string
		args    args
		wantLen int
		wantErr bool
	}{
		{
			name: "get by preset Id",
			args: args{
				conds: []Condition{
					ById(1),
				},
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "get by preset UId",
			args: args{
				conds: []Condition{
					ByUid("0195fad1-6bd6-765a-8a25-6be0bd03e9ce"),
				},
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "get by wrong preset Id",
			args: args{
				conds: []Condition{
					ById(0),
				},
			},
			wantLen: 0,
			wantErr: false,
		},
		{
			name:    "get without id",
			wantLen: 0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := d.GetPreset(tt.args.conds...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Sqlite.GetPreset() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantLen != 0 && got == nil {
				t.Errorf("Sqlite.ListInstruments() got nil, want len = %v", tt.wantLen)
			}
			if tt.wantLen == 0 && got != nil {
				t.Errorf("Sqlite.ListInstruments() got data, want len = %v", tt.wantLen)
			}
		})
	}
}

func TestSqlite_ListPresets(t *testing.T) {
	d, err := NewSqlite(getDBPath())
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	defer d.Close()

	type args struct {
		conds []Condition
	}
	tests := []struct {
		name    string
		args    args
		wantLen int
		wantErr bool
	}{
		{
			name:    "without conditions",
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "by preset id",
			args: args{
				conds: []Condition{
					ById(1),
				},
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "by kit id",
			args: args{
				conds: []Condition{
					ByKitId(1),
				},
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "by wrong kit id",
			args: args{
				conds: []Condition{
					ByKitId(0),
				},
			},
			wantLen: 0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := d.ListPresets(tt.args.conds...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Sqlite.ListPresets() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantLen != 0 && (got == nil || len(*got) == 0) {
				t.Errorf("Sqlite.ListPresets() got nil, want len = %v", tt.wantLen)
			}
			if tt.wantLen == 0 && got != nil && len(*got) > 0 {
				t.Errorf("Sqlite.ListPresets() got data, want len = %v", tt.wantLen)
			}
		})
	}
}
