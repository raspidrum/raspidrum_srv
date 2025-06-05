package preset

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	m "github.com/raspidrum-srv/internal/model"
)

func Test_PrepareToLoad(t *testing.T) {
	type args struct {
		preset   *m.KitPreset
		mididevs []m.MIDIDevice
	}
	type testCase struct {
		name             string
		testData         string
		args             args
		want             m.KitPreset
		wantErr          bool
		expectedControls m.ExpectedControls
	}
	tests := []testCase{
		{
			name:     "channel with one instrument without layers",
			testData: "single_instrument.yaml",
			args: args{
				mididevs: []m.MIDIDevice{
					&MockMMIDIDevice{},
				},
			},
			want: m.KitPreset{
				Uid:  "preset-1",
				Name: "Single Instrument",
				Channels: []m.PresetChannel{
					{
						Key:  "ch1",
						Name: "Kick",
						Controls: map[string]*m.PresetControl{
							"volume": {Key: "c0volume", Name: "Volume", Type: "volume", Value: 1.00},
							"pan":    {Key: "c0pan", Name: "Pan", Type: "pan", Value: 0.00},
						},
					},
				},
				Instruments: []m.PresetInstrument{
					{
						Name:       "Kick",
						ChannelKey: "ch1",
						MidiKey:    "kick1",
						MidiNote:   36,
						Controls: map[string]*m.PresetControl{
							"volume": {Key: "i0volume", Name: "Volume", MidiCC: 30, CfgKey: "KICKV", Type: "volume", Value: 95},
							"pan":    {Key: "i0pan", Name: "Pan", MidiCC: 10, CfgKey: "KICKP", Type: "pan", Value: 54},
						},
					},
				},
			},
			wantErr: false,
			expectedControls: m.ExpectedControls{
				"c0volume": {Key: "c0volume", Name: "Volume", Type: "volume", Value: 1.00},
				"c0pan":    {Key: "c0pan", Name: "Pan", Type: "pan", Value: 0.00},
				"i0volume": {Key: "i0volume", Name: "Volume", Type: "volume", MidiCC: 30, CfgKey: "KICKV", Value: 95},
				"i0pan":    {Key: "i0pan", Name: "Pan", Type: "pan", MidiCC: 10, CfgKey: "KICKP", Value: 54},
			},
		},
		{
			name:     "with layers",
			testData: "single instr_with_layers.yaml",
			args: args{
				mididevs: []m.MIDIDevice{
					&MockMMIDIDevice{},
				},
			},
			want: m.KitPreset{
				Channels: []m.PresetChannel{
					{
						Key: "ch1",
						Controls: map[string]*m.PresetControl{
							"volume": {Key: "c0volume", Type: "volume", Value: 66},
							"pan":    {Key: "c0pan", Type: "pan", Value: 0.00},
						},
					},
				},
				Instruments: []m.PresetInstrument{
					{
						Name:       "Ride",
						ChannelKey: "ch1",
						Controls: map[string]*m.PresetControl{
							"pan":    {Key: "i0pan", MidiCC: 105, CfgKey: "RI17P", Type: "pan", Value: 75},
							"pitch":  {Key: "i0pitch", MidiCC: 16, CfgKey: "RI17T", Type: "pitch", Value: 120},
							"volume": {Key: "i0volume", Type: "volume", Value: 111},
						},
						Layers: map[string]m.PresetLayer{
							"bell": {
								MidiKey:    "ride1_bell",
								CfgMidiKey: "RI17BKEY",
								MidiNote:   53,
								Controls: map[string]*m.PresetControl{
									"volume": {Key: "i0l0volume", MidiCC: 104, CfgKey: "RI17BV", Type: "volume", Value: 80},
								},
							},
							"edge": {
								MidiKey:    "ride1_edge",
								CfgMidiKey: "RI17EKEY",
								MidiNote:   51,
								Controls: map[string]*m.PresetControl{
									"volume": {Key: "i0l1volume", MidiCC: 103, CfgKey: "RI17EV", Type: "volume", Value: 90},
								},
							},
						},
					},
				},
			},
			wantErr: false,
			expectedControls: m.ExpectedControls{
				"c0volume":   {Key: "c0volume", Type: "volume", Value: 66},
				"c0pan":      {Key: "c0pan", Type: "pan", Value: 0.00},
				"i0pan":      {Key: "i0pan", MidiCC: 105, CfgKey: "RI17P", Type: "pan", Value: 75},
				"i0pitch":    {Key: "i0pitch", MidiCC: 16, CfgKey: "RI17T", Type: "pitch", Value: 120},
				"i0volume":   {Key: "i0volume", Type: "volume", Value: 111},
				"i0l0volume": {Key: "i0l0volume", MidiCC: 104, CfgKey: "RI17BV", Type: "volume", Value: 80},
				"i0l1volume": {Key: "i0l1volume", MidiCC: 103, CfgKey: "RI17EV", Type: "volume", Value: 90},
			},
		},
		{
			name:     "two instruments",
			testData: "two_instruments.yaml",
			args: args{
				mididevs: []m.MIDIDevice{
					&MockMMIDIDevice{},
				},
			},
			want: m.KitPreset{
				Uid:  "preset-2",
				Name: "multiple instruments in channel",
				Channels: []m.PresetChannel{
					{
						Key:  "ch1",
						Name: "Drums",
						Controls: map[string]*m.PresetControl{
							"volume": {Key: "c0volume", Name: "Volume", Type: "volume", Value: 1.00},
							"pan":    {Key: "c0pan", Name: "Pan", Type: "pan", Value: 0.00},
						},
					},
				},
				Instruments: []m.PresetInstrument{
					{
						Name:       "Kick",
						ChannelKey: "ch1",
						MidiKey:    "kick1",
						MidiNote:   36,
						Controls: map[string]*m.PresetControl{
							"volume": {Key: "i0volume", Name: "Volume", MidiCC: 30, CfgKey: "KICKV", Type: "volume", Value: 95},
							"pan":    {Key: "i0pan", Name: "Pan", MidiCC: 10, CfgKey: "KICKP", Type: "pan", Value: 54},
						},
					},
					{
						Name:       "Tom",
						ChannelKey: "ch1",
						MidiKey:    "tom1",
						MidiNote:   48,
						Controls: map[string]*m.PresetControl{
							"volume": {Key: "i1volume", Name: "Volume", MidiCC: 31, CfgKey: "TOM1V", Type: "volume", Value: 87},
						},
					},
				},
			},
			wantErr: false,
			expectedControls: m.ExpectedControls{
				"c0volume": {Key: "c0volume", Name: "Volume", Type: "volume", Value: 1.00},
				"c0pan":    {Key: "c0pan", Name: "Pan", Type: "pan", Value: 0.00},
				"i0volume": {Key: "i0volume", Name: "Volume", MidiCC: 30, CfgKey: "KICKV", Type: "volume", Value: 95},
				"i0pan":    {Key: "i0pan", Name: "Pan", MidiCC: 10, CfgKey: "KICKP", Type: "pan", Value: 54},
				"i1volume": {Key: "i1volume", Name: "Volume", MidiCC: 31, CfgKey: "TOM1V", Type: "volume", Value: 87},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			preset := loadPresetFromYAML(t, tt.testData)
			tt.args.preset = preset
			if err := preset.PrepareToLoad(tt.args.mididevs); (err != nil) != tt.wantErr {
				t.Errorf("PrepareToLoad() error = %v, wantErr %v", err, tt.wantErr)
			}
			if diff := cmp.Diff(tt.want, *preset,
				cmpopts.IgnoreUnexported(m.KitPreset{}),
				cmpopts.IgnoreUnexported(m.PresetControl{}),
				cmpopts.IgnoreUnexported(m.PresetInstrument{}),
				cmpopts.IgnoreUnexported(m.PresetChannel{}),
				cmpopts.IgnoreUnexported(m.InstrumentRef{}),
				cmpopts.IgnoreUnexported(m.Layer{}),
				cmpopts.IgnoreUnexported(m.Control{}),
				cmpopts.IgnoreFields(m.PresetInstrument{}, "Instrument"),
				cmpopts.IgnoreFields(m.PresetInstrument{}, "Id"),
			); diff != "" {
				t.Errorf("PrepareToLoad() mismatch (-want +got):\n%s", diff)
			}

			// Verify controls state using the method from KitPreset
			if diff := m.VerifyControlsForTest(preset, tt.expectedControls); diff != "" {
				t.Error("Controls state does not match expected:", diff)
			}
		})
	}
}
