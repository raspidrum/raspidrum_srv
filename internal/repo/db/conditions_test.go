package db

import (
	"reflect"
	"testing"
)

func TestEq(t *testing.T) {
	type args struct {
		field  string
		inargs []interface{}
	}
	tests := []struct {
		name     string
		args     args
		wantSql  string
		wantArgs []interface{}
		wantErr  bool
	}{
		{
			name: "simple",
			args: args{
				field:  "id",
				inargs: append([]any{}, 1),
			},
			wantSql:  "id = ?",
			wantArgs: append([]any{}, 1),
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cond := Eq(tt.args.field, tt.args.inargs...)
			gotSql, gotArgs, err := cond()
			if (err != nil) != tt.wantErr {
				t.Errorf("Eq() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotSql != tt.wantSql {
				t.Errorf("Eq() gotSql = %v, want %v", gotSql, tt.wantSql)
			}
			if !reflect.DeepEqual(gotArgs, tt.wantArgs) {
				t.Errorf("Eq() gotArgs = %v, want %v", gotArgs, tt.wantArgs)
			}
		})
	}
}

func TestIn(t *testing.T) {
	type args struct {
		field  string
		inargs []interface{}
	}
	tests := []struct {
		name     string
		args     args
		wantSql  string
		wantArgs []interface{}
	}{
		{
			name: "one arg",
			args: args{
				field:  "tag",
				inargs: append([]any{}, "rock"),
			},
			wantSql:  "tag in (?)",
			wantArgs: append([]any{}, "rock"),
		},
		{
			name: "two arg",
			args: args{
				field:  "tag",
				inargs: append([]any{}, "rock", "retro"),
			},
			wantSql:  "tag in (?, ?)",
			wantArgs: append([]any{}, "rock", "retro"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cond := In(tt.args.field, tt.args.inargs...)
			gotSql, gotArgs, err := cond()
			if err != nil {
				t.Errorf("In() error = %v", err)
				return
			}
			if gotSql != tt.wantSql {
				t.Errorf("In() gotSql = %v, want %v", gotSql, tt.wantSql)
			}
			if !reflect.DeepEqual(gotArgs, tt.wantArgs) {
				t.Errorf("In() gotArgs = %v, want %v", gotArgs, tt.wantArgs)
			}
		})
	}
}

func Test_buildConditions(t *testing.T) {
	type args struct {
		conds []Condition
	}
	tests := []struct {
		name     string
		args     args
		wantSql  string
		wantArgs []interface{}
		wantErr  bool
	}{
		{
			name: "one arg",
			args: args{
				conds: []Condition{
					Eq("id", 1),
				},
			},
			wantSql:  "where id = ?",
			wantArgs: []any{1},
			wantErr:  false,
		},
		{
			name: "two args",
			args: args{
				conds: []Condition{
					Eq("id", 1),
					Eq("kit", 2),
				},
			},
			wantSql:  "where id = ? and kit = ?",
			wantArgs: []any{1, 2},
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSql, gotArgs, err := buildConditions(tt.args.conds...)
			if (err != nil) != tt.wantErr {
				t.Errorf("buildConditions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotSql != tt.wantSql {
				t.Errorf("buildConditions() gotSql = %v, want %v", gotSql, tt.wantSql)
			}
			if !reflect.DeepEqual(gotArgs, tt.wantArgs) {
				t.Errorf("buildConditions() gotArgs = %v, want %v", gotArgs, tt.wantArgs)
			}
		})
	}
}
