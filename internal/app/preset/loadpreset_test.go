package preset

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/goccy/go-yaml"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	m "github.com/raspidrum-srv/internal/model"
)

type MockMMIDIDevice struct{}

func (m *MockMMIDIDevice) Name() string {
	return "Dummy"
}

func (m *MockMMIDIDevice) DevID() string {
	return "0:0"
}

func (m *MockMMIDIDevice) GetKeysMapping() (map[string]int, error) {
	return map[string]int{
		"kick1":      36,
		"snare":      38,
		"ride1_edge": 51,
		"ride1_bell": 53,
		"tom1":       48,
	}, nil
}

func loadPresetFromYAML(t *testing.T, filename string) *m.KitPreset {
	t.Helper()
	absPath, err := filepath.Abs("../../../testdata/preset_load/" + filename)
	if err != nil {
		t.Fatalf("Failed to get absolute path: %v", err)
	}
	data, err := os.ReadFile(absPath)
	if err != nil {
		t.Fatalf("Failed to read test data file: %v", err)
	}
	var preset m.KitPreset
	if err := yaml.Unmarshal(data, &preset); err != nil {
		t.Fatalf("Failed to unmarshal YAML: %v", err)
	}
	return &preset
}

func Test_PrepareToLoad(t *testing.T) {
	type args struct {
		pst      *m.KitPreset
		mididevs []m.MIDIDevice
	}
	type testCase struct {
		name             string
		yamlFile         string
		args             args
		want             m.KitPreset
		wantErr          bool
		expectedControls map[string]struct {
			Key    string
			Owner  m.ControlOwner
			MidiCC int
			CfgKey string
			Type   string
			Value  float32
		}
	}
	tests := []testCase{
		{
			name:     "channel with one instrument without layers",
			yamlFile: "single_instrument.yaml",
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
							"volume": {Key: "i0volume", MidiCC: 30, CfgKey: "KICKV", Type: "volume", Value: 95},
						},
					},
				},
			},
			wantErr: false,
			expectedControls: map[string]struct {
				Key    string
				Owner  m.ControlOwner
				MidiCC int
				CfgKey string
				Type   string
				Value  float32
			}{
				"c0volume": {Key: "c0volume", Type: "volume", Value: 66},
				"i0volume": {Key: "i0volume", Type: "volume", MidiCC: 30, CfgKey: "KICKV", Value: 95},
			},
		},
		{
			name:     "with layers",
			yamlFile: "single instr_with_layers.yaml",
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
			expectedControls: map[string]struct {
				Key    string
				Owner  m.ControlOwner
				MidiCC int
				CfgKey string
				Type   string
				Value  float32
			}{
				"c0volume":   {Key: "c0volume", Type: "volume", Value: 66},
				"i0pan":      {Key: "i0pan", MidiCC: 105, CfgKey: "RI17P", Type: "pan", Value: 75},
				"i0pitch":    {Key: "i0pitch", MidiCC: 16, CfgKey: "RI17T", Type: "pitch", Value: 120},
				"i0volume":   {Key: "i0volume", Type: "volume", Value: 111},
				"i0l0volume": {Key: "i0l0volume", MidiCC: 104, CfgKey: "RI17BV", Type: "volume", Value: 80},
				"i0l1volume": {Key: "i0l1volume", MidiCC: 103, CfgKey: "RI17EV", Type: "volume", Value: 90},
			},
		},
		{
			name:     "two instruments",
			yamlFile: "two_instruments.yaml",
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
							"volume": {Key: "i0volume", MidiCC: 30, CfgKey: "KICKV", Type: "volume", Value: 95},
						},
					},
					{
						Name:       "Tom",
						ChannelKey: "ch1",
						MidiKey:    "tom1",
						MidiNote:   48,
						Controls: map[string]*m.PresetControl{
							"volume": {Key: "i1volume", MidiCC: 31, CfgKey: "TOM1V", Type: "volume", Value: 87},
						},
					},
				},
			},
			wantErr: false,
			expectedControls: map[string]struct {
				Key    string
				Owner  m.ControlOwner
				MidiCC int
				CfgKey string
				Type   string
				Value  float32
			}{
				"c0volume": {Key: "c0volume", Type: "volume", Value: 66},
				"i0volume": {Key: "i0volume", MidiCC: 30, CfgKey: "KICKV", Type: "volume", Value: 95},
				"i1volume": {Key: "i1volume", MidiCC: 31, CfgKey: "TOM1V", Type: "volume", Value: 87},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			preset := loadPresetFromYAML(t, tt.yamlFile)
			tt.args.pst = preset
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
				cmpopts.IgnoreFields(m.PresetInstrument{}, "Instrument")); diff != "" {
				t.Errorf("PrepareToLoad() mismatch (-want +got):\n%s", diff)
			}

			// Verify controls state using the method from KitPreset
			// TODO: collect expected controls from preset and compare with actual controls
			if diff := m.VerifyControlsForTest(preset, tt.expectedControls); diff != "" {
				t.Error("Controls state does not match expected:", diff)
			}
		})
	}
}
