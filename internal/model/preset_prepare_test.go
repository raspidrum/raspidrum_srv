package model

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func Test_PrepareToLoad(t *testing.T) {
	type args struct {
		preset   *KitPreset
		mididevs []MIDIDevice
	}
	type testCase struct {
		name             string
		testData         string
		args             args
		want             KitPreset
		wantErr          bool
		expectedControls ExpectedControls
	}
	tests := []testCase{
		{
			name:     "channel with one instrument without layers",
			testData: "single_instrument.yaml",
			args: args{
				mididevs: []MIDIDevice{
					&MockMMIDIDevice{},
				},
			},
			want: KitPreset{
				Uid:  "preset-1",
				Name: "Single Instrument",
				Channels: []PresetChannel{
					{
						Key:  "ch1",
						Name: "Kick",
						Controls: map[string]*PresetControl{
							"volume": {Key: "c0volume", Name: "Volume", Type: "volume", Value: 1.00},
							"pan":    {Key: "c0pan", Name: "Pan", Type: "pan", Value: 0.00},
						},
					},
				},
				Instruments: []PresetInstrument{
					{
						Name:       "Kick",
						ChannelKey: "ch1",
						MidiKey:    "kick1",
						MidiNote:   36,
						Controls: map[string]*PresetControl{
							"volume": {Key: "i0volume", Name: "Volume", MidiCC: 30, CfgKey: "KICKV", Type: "volume", Value: 95},
							"pan":    {Key: "i0pan", Name: "Pan", MidiCC: 10, CfgKey: "KICKP", Type: "pan", Value: 54},
						},
					},
				},
			},
			wantErr: false,
			expectedControls: ExpectedControls{
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
				mididevs: []MIDIDevice{
					&MockMMIDIDevice{},
				},
			},
			want: KitPreset{
				Channels: []PresetChannel{
					{
						Key: "ch1",
						Controls: map[string]*PresetControl{
							"volume": {Key: "c0volume", Type: "volume", Value: 66},
							"pan":    {Key: "c0pan", Type: "pan", Value: 0.00},
						},
					},
				},
				Instruments: []PresetInstrument{
					{
						Name:       "Ride",
						ChannelKey: "ch1",
						Controls: map[string]*PresetControl{
							"pan":    {Key: "i0pan", MidiCC: 105, CfgKey: "RI17P", Type: "pan", Value: 75},
							"pitch":  {Key: "i0pitch", MidiCC: 16, CfgKey: "RI17T", Type: "pitch", Value: 120},
							"volume": {Key: "i0volume", Type: "volume", Value: 111},
						},
						Layers: map[string]PresetLayer{
							"bell": {
								MidiKey:    "ride1_bell",
								CfgMidiKey: "RI17BKEY",
								MidiNote:   53,
								Controls: map[string]*PresetControl{
									"volume": {Key: "i0l0volume", MidiCC: 104, CfgKey: "RI17BV", Type: "volume", Value: 80},
								},
							},
							"edge": {
								MidiKey:    "ride1_edge",
								CfgMidiKey: "RI17EKEY",
								MidiNote:   51,
								Controls: map[string]*PresetControl{
									"volume": {Key: "i0l1volume", MidiCC: 103, CfgKey: "RI17EV", Type: "volume", Value: 90},
								},
							},
						},
					},
				},
			},
			wantErr: false,
			expectedControls: ExpectedControls{
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
				mididevs: []MIDIDevice{
					&MockMMIDIDevice{},
				},
			},
			want: KitPreset{
				Uid:  "preset-2",
				Name: "multiple instruments in channel",
				Channels: []PresetChannel{
					{
						Key:  "ch1",
						Name: "Drums",
						Controls: map[string]*PresetControl{
							"volume": {Key: "c0volume", Name: "Volume", Type: "volume", Value: 1.00},
							"pan":    {Key: "c0pan", Name: "Pan", Type: "pan", Value: 0.00},
						},
					},
				},
				Instruments: []PresetInstrument{
					{
						Name:       "Kick",
						ChannelKey: "ch1",
						MidiKey:    "kick1",
						MidiNote:   36,
						Controls: map[string]*PresetControl{
							"volume": {Key: "i0volume", Name: "Volume", MidiCC: 30, CfgKey: "KICKV", Type: "volume", Value: 95},
							"pan":    {Key: "i0pan", Name: "Pan", MidiCC: 10, CfgKey: "KICKP", Type: "pan", Value: 54},
						},
					},
					{
						Name:       "Tom",
						ChannelKey: "ch1",
						MidiKey:    "tom1",
						MidiNote:   48,
						Controls: map[string]*PresetControl{
							"volume": {Key: "i1volume", Name: "Volume", MidiCC: 31, CfgKey: "TOM1V", Type: "volume", Value: 87},
						},
					},
				},
			},
			wantErr: false,
			expectedControls: ExpectedControls{
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
				cmpopts.IgnoreUnexported(KitPreset{}),
				cmpopts.IgnoreUnexported(PresetControl{}),
				cmpopts.IgnoreUnexported(PresetInstrument{}),
				cmpopts.IgnoreUnexported(PresetChannel{}),
				cmpopts.IgnoreUnexported(InstrumentRef{}),
				cmpopts.IgnoreUnexported(Layer{}),
				cmpopts.IgnoreUnexported(Control{}),
				cmpopts.IgnoreFields(PresetInstrument{}, "Instrument"),
				cmpopts.IgnoreFields(PresetInstrument{}, "Id"),
			); diff != "" {
				t.Errorf("PrepareToLoad() mismatch (-want +got):\n%s", diff)
			}

			// Verify controls state using the method from KitPreset
			// TODO: collect expected controls from preset and compare with actual controls
			if diff := VerifyControlsForTest(preset, tt.expectedControls); diff != "" {
				t.Error("Controls state does not match expected:", diff)
			}
		})
	}
}
