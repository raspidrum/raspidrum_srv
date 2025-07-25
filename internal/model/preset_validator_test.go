package model

import (
	"testing"
)

func TestPresetControl_ValidateType(t *testing.T) {
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

func TestPresetControl_ValidateValue(t *testing.T) {
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
			name:    "volume with midiCC",
			fields:  fields{Type: "volume", MidiCC: 123, Value: 20},
			wantErr: false,
		},
		{
			name:    "volume without midiCC",
			fields:  fields{Type: "volume", Value: 0.2},
			wantErr: false,
		},
		{
			name:    "pan with midiCC",
			fields:  fields{Type: "pan", MidiCC: 123, Value: 10},
			wantErr: false,
		},
		{
			name:    "pan without midiCC",
			fields:  fields{Type: "pan", Value: 0.5},
			wantErr: false,
		},
		{
			name:    "volume with midiCC, value over of range",
			fields:  fields{Type: "volume", MidiCC: 123, Value: 207},
			wantErr: true,
		},
		{
			name:    "volume with midiCC, value lower of range",
			fields:  fields{Type: "volume", MidiCC: 123, Value: -10},
			wantErr: true,
		},
		{
			name:    "pan with midiCC, value lower of range",
			fields:  fields{Type: "pan", MidiCC: 123, Value: -10},
			wantErr: true,
		},
		{
			name:    "pan without midiCC, value lower of range",
			fields:  fields{Type: "pan", Value: 1.5},
			wantErr: true,
		},
		{
			name:    "volume without midiCC, value lower of range",
			fields:  fields{Type: "volume", Value: -0.5},
			wantErr: true,
		},
		{
			name:    "volume without midiCC, value over of range",
			fields:  fields{Type: "volume", Value: 1.5},
			wantErr: true,
		},
		{
			name:    "pan without midiCC, value over of range",
			fields:  fields{Type: "pan", Value: 1.5},
			wantErr: true,
		},
		{
			name:    "other control with midiCC, value lower of range",
			fields:  fields{Type: "other", MidiCC: 123, Value: -10},
			wantErr: true,
		},
		{
			name:    "other control with midiCC, value over of range",
			fields:  fields{Type: "other", MidiCC: 123, Value: 207},
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

func TestPresetLayer_Validate(t *testing.T) {
	type fields struct {
		Name       string
		MidiKey    string
		CfgMidiKey string
		MidiNote   int
		Controls   map[string]*PresetControl
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "vol and pan with midiCC",
			fields: fields{
				Controls: map[string]*PresetControl{
					"volume": {Type: "volume", MidiCC: 123},
					"pan":    {Type: "pan", MidiCC: 134},
				},
			},
			wantErr: false,
		},
		{
			name: "vol with midiCC, no pan",
			fields: fields{
				Controls: map[string]*PresetControl{
					"volume": {Type: "volume", MidiCC: 123},
				},
			},
			wantErr: false,
		},
		{
			name: "vol with midiCC, pan without midiCC",
			fields: fields{
				Controls: map[string]*PresetControl{
					"volume": {Type: "volume", MidiCC: 123},
					"pan":    {Type: "pan"},
				},
			},
			wantErr: true,
		},
		{
			name: "no vol",
			fields: fields{
				Controls: map[string]*PresetControl{},
			},
			wantErr: true,
		},
		{
			name: "vol without midiCC",
			fields: fields{
				Controls: map[string]*PresetControl{
					"volume": {Type: "volume"},
				},
			},
			wantErr: true,
		},
		{
			name: "other control without midiCC",
			fields: fields{
				Controls: map[string]*PresetControl{
					"pitch": {Type: "pitch"},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PresetLayer{
				Name:       tt.fields.Name,
				MidiKey:    tt.fields.MidiKey,
				CfgMidiKey: tt.fields.CfgMidiKey,
				MidiNote:   tt.fields.MidiNote,
				Controls:   tt.fields.Controls,
			}
			if err := p.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("PresetLayer.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestKitPreset_Validate(t *testing.T) {
	type fields struct {
		Uid         string
		Kit         KitRef
		Name        string
		Channels    []PresetChannel
		Instruments []PresetInstrument
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "one instr in channel, instr without layers. Instr with vol with MidiCC",
			fields: fields{
				Channels: []PresetChannel{{Key: "1"}},
				Instruments: []PresetInstrument{
					{ChannelKey: "1",
						Controls: map[string]*PresetControl{
							"volume": {Type: "volume", MidiCC: 123},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "one instr in channel, instr without layers. Instr with vol without MidiCC",
			fields: fields{
				Channels: []PresetChannel{{Key: "1"}},
				Instruments: []PresetInstrument{
					{ChannelKey: "1",
						Controls: map[string]*PresetControl{
							"volume": {Type: "volume"},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "one instr in channel, instr without layers. Instr without vol",
			fields: fields{
				Channels: []PresetChannel{{Key: "1"}},
				Instruments: []PresetInstrument{
					{ChannelKey: "1"},
				},
			},
			wantErr: false,
		},
		{
			name: "one instr in channel, instr with layers. Instr with vol with MidiCC",
			fields: fields{
				Channels: []PresetChannel{{Key: "1"}},
				Instruments: []PresetInstrument{
					{ChannelKey: "1",
						Controls: map[string]*PresetControl{
							"volume": {Type: "volume", MidiCC: 123},
						},
						Layers: map[string]PresetLayer{
							"top": {
								Controls: map[string]*PresetControl{
									"volume": {Type: "volume", MidiCC: 123},
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "one instr in channel, instr with layers. Instr with vol without MidiCC",
			fields: fields{
				Channels: []PresetChannel{{Key: "1"}},
				Instruments: []PresetInstrument{
					{ChannelKey: "1",
						Controls: map[string]*PresetControl{
							"volume": {Type: "volume"},
						},
						Layers: map[string]PresetLayer{
							"top": {
								Controls: map[string]*PresetControl{
									"volume": {Type: "volume", MidiCC: 123},
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "one instr in channel, instr with layers. Instr without vol",
			fields: fields{
				Channels: []PresetChannel{{Key: "1"}},
				Instruments: []PresetInstrument{
					{ChannelKey: "1",
						Controls: map[string]*PresetControl{},
						Layers: map[string]PresetLayer{
							"top": {
								Controls: map[string]*PresetControl{
									"volume": {Type: "volume", MidiCC: 123},
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "many instr in channel, instr without layers. Instr with vol with MidiCC",
			fields: fields{
				Channels: []PresetChannel{{Key: "1"}},
				Instruments: []PresetInstrument{
					{ChannelKey: "1",
						Name: "tom1",
						Controls: map[string]*PresetControl{
							"volume": {Type: "volume", MidiCC: 123},
							"pan":    {Type: "pan", MidiCC: 123},
						},
					},
					{ChannelKey: "1",
						Name: "tom2",
						Controls: map[string]*PresetControl{
							"volume": {Type: "volume", MidiCC: 123},
							"pan":    {Type: "pan", MidiCC: 123},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "many instr in channel, instr without layers. Instr with vol without MidiCC",
			fields: fields{
				Channels: []PresetChannel{{Key: "1"}},
				Instruments: []PresetInstrument{
					{ChannelKey: "1",
						Name: "tom1",
						Controls: map[string]*PresetControl{
							"volume": {Type: "volume"},
							"pan":    {Type: "pan", MidiCC: 123},
						},
					},
					{ChannelKey: "1",
						Name: "tom2",
						Controls: map[string]*PresetControl{
							"volume": {Type: "volume", MidiCC: 123},
							"pan":    {Type: "pan"},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "many instr in channel, instr without layers. Instr without vol",
			fields: fields{
				Channels: []PresetChannel{{Key: "1"}},
				Instruments: []PresetInstrument{
					{ChannelKey: "1",
						Name: "tom1",
						Controls: map[string]*PresetControl{
							"pan": {Type: "pan", MidiCC: 123},
						},
					},
					{ChannelKey: "1",
						Name: "tom2",
						Controls: map[string]*PresetControl{
							"volume": {Type: "volume", MidiCC: 123},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "many instr in channel, instr with layers. Instr with vol with MidiCC",
			fields: fields{
				Channels: []PresetChannel{{Key: "1"}},
				Instruments: []PresetInstrument{
					{ChannelKey: "1",
						Name: "tom1",
						Controls: map[string]*PresetControl{
							"volume": {Type: "volume", MidiCC: 123},
							"pan":    {Type: "pan", MidiCC: 123},
						},
						Layers: map[string]PresetLayer{
							"top": {
								Controls: map[string]*PresetControl{
									"volume": {Type: "volume", MidiCC: 123},
								},
							},
						},
					},
					{ChannelKey: "1",
						Name: "tom2",
						Controls: map[string]*PresetControl{
							"volume": {Type: "volume", MidiCC: 123},
							"pan":    {Type: "pan", MidiCC: 123},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "many instr in channel, instr with layers. Instr with vol without MidiCC",
			fields: fields{
				Channels: []PresetChannel{{Key: "1"}},
				Instruments: []PresetInstrument{
					{ChannelKey: "1",
						Name: "tom1",
						Controls: map[string]*PresetControl{
							"volume": {Type: "volume"},
							"pan":    {Type: "pan", MidiCC: 123},
						},
						Layers: map[string]PresetLayer{
							"top": {
								Controls: map[string]*PresetControl{
									"volume": {Type: "volume", MidiCC: 123},
								},
							},
						},
					},
					{ChannelKey: "1",
						Name: "tom2",
						Controls: map[string]*PresetControl{
							"volume": {Type: "volume", MidiCC: 123},
							"pan":    {Type: "pan", MidiCC: 123},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "many instr in channel, instr with layers. Instr without vol",
			fields: fields{
				Channels: []PresetChannel{{Key: "1"}},
				Instruments: []PresetInstrument{
					{ChannelKey: "1",
						Name: "tom1",
						Controls: map[string]*PresetControl{
							"pan": {Type: "pan", MidiCC: 123},
						},
						Layers: map[string]PresetLayer{
							"top": {
								Controls: map[string]*PresetControl{
									"volume": {Type: "volume", MidiCC: 123},
								},
							},
						},
					},
					{ChannelKey: "1",
						Name: "tom2",
						Controls: map[string]*PresetControl{
							"volume": {Type: "volume", MidiCC: 123},
							"pan":    {Type: "pan", MidiCC: 123},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "other control must have midiCC",
			fields: fields{
				Channels: []PresetChannel{{Key: "1"}},
				Instruments: []PresetInstrument{
					{ChannelKey: "1",
						Controls: map[string]*PresetControl{
							"pitch": {Type: "pitch"},
						},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &KitPreset{
				Uid:         tt.fields.Uid,
				Kit:         tt.fields.Kit,
				Name:        tt.fields.Name,
				Channels:    tt.fields.Channels,
				Instruments: tt.fields.Instruments,
			}
			if err := p.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("KitPreset.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
