package model

import (
	"reflect"
	"testing"
)

func TestKitPreset_IndexInstruments(t *testing.T) {
	type fields struct {
		Channels    []PresetChannel
		Instruments []PresetInstrument
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "correct ChannelKeys",
			fields: fields{
				Channels: []PresetChannel{
					{Key: "a"}, {Key: "b"},
				},
				Instruments: []PresetInstrument{
					{Name: "A", ChannelKey: "a"},
					{Name: "B", ChannelKey: "b"},
				},
			},
			wantErr: false,
		},
		{
			name: "incorrect ChannelKeys",
			fields: fields{
				Channels: []PresetChannel{
					{Key: "a"}, {Key: "b"},
				},
				Instruments: []PresetInstrument{
					{Name: "A", ChannelKey: "1"},
					{Name: "B", ChannelKey: "b"},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &KitPreset{
				Channels:    tt.fields.Channels,
				Instruments: tt.fields.Instruments,
			}
			if err := p.indexInstruments(); (err != nil) != tt.wantErr {
				t.Errorf("KitPreset.IndexInstruments() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestKitPreset_GetChannelInstrumentsByKey(t *testing.T) {
	type fields struct {
		Channels    []PresetChannel
		Instruments []PresetInstrument
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*PresetInstrument
		wantErr bool
	}{
		{
			name: "correct Channel key",
			fields: fields{
				Channels: []PresetChannel{
					{Key: "a"}, {Key: "b"},
				},
				Instruments: []PresetInstrument{
					{Name: "A", ChannelKey: "a"},
					{Name: "B", ChannelKey: "b"},
				},
			},
			args:    args{key: "b"},
			want:    []*PresetInstrument{},
			wantErr: false,
		},
		{
			name: "incorrect Channel key",
			fields: fields{
				Channels: []PresetChannel{
					{Key: "a"}, {Key: "b"},
				},
				Instruments: []PresetInstrument{
					{Name: "A", ChannelKey: "a"},
					{Name: "B", ChannelKey: "b"},
				},
			},
			args:    args{key: "c"},
			want:    []*PresetInstrument{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &KitPreset{
				Channels:    tt.fields.Channels,
				Instruments: tt.fields.Instruments,
			}
			if !tt.wantErr {
				for i, v := range tt.fields.Instruments {
					if v.ChannelKey == tt.args.key {
						tt.want = append(tt.want, &tt.fields.Instruments[i])
					}
				}
			}

			got, err := p.GetChannelInstrumentsByKey(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("KitPreset.GetChannelInstrumentsByKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("KitPreset.GetChannelInstrumentsByKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
