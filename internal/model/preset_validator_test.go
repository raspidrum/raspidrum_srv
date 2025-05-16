package model

import "testing"

func TestPresetControl_Validate(t *testing.T) {
	type fields struct {
		Name   string
		Type   string
		MidiCC int
		CfgKey string
		Value  float32
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "match",
			fields:  fields{Type: "volume"},
			wantErr: false,
		},
		{
			name:    "doesn't match",
			fields:  fields{Type: "level"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &PresetControl{
				Name:   tt.fields.Name,
				Type:   tt.fields.Type,
				MidiCC: tt.fields.MidiCC,
				CfgKey: tt.fields.CfgKey,
				Value:  tt.fields.Value,
			}
			if err := c.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("PresetControl.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
