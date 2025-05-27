package db

import (
	"reflect"
	"testing"
)

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
