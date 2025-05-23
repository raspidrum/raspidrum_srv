package preset

import (
	"path"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	m "github.com/raspidrum-srv/internal/model"
	"github.com/raspidrum-srv/internal/repo"
	"github.com/raspidrum-srv/internal/repo/db"
	"github.com/spf13/afero"
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
	}, nil
}

func Test_augmentFromInstrument(t *testing.T) {
	type args struct {
		pst      *m.KitPreset
		mididevs []m.MIDIDevice
	}
	type testCase struct {
		name            string
		args            args
		want            m.KitPreset
		wantErr         bool
		expectedControl map[string]struct {
			Owner  m.ControlOwner
			MidiCC int
			CfgKey string
			Type   string
			Value  float32
		}
	}
	tests := []testCase{
		{
			name: "wo layers",
			args: args{
				pst: &m.KitPreset{
					Instruments: []m.PresetInstrument{
						{
							Instrument: m.InstrumentRef{
								CfgMidiKey: "KEYKICK",
								Controls: map[string]m.Control{
									"volume": {CfgKey: "KICKV"},
								},
							},
							MidiKey: "kick1",
							Controls: map[string]m.PresetControl{
								"volume": {
									MidiCC: 30,
									Type:   "volume",
								},
							},
						},
					},
				},
				mididevs: []m.MIDIDevice{
					&MockMMIDIDevice{},
				},
			},
			want: m.KitPreset{
				Instruments: []m.PresetInstrument{
					{
						Instrument: m.InstrumentRef{
							CfgMidiKey: "KEYKICK",
							Controls: map[string]m.Control{
								"volume": {CfgKey: "KICKV"},
							},
						},
						MidiKey:  "kick1",
						MidiNote: 36,
						Controls: map[string]m.PresetControl{
							"volume": {
								MidiCC: 30,
								CfgKey: "KICKV",
								Type:   "volume",
							},
						},
					},
				},
			},
			wantErr: false,
			expectedControl: map[string]struct {
				Owner  m.ControlOwner
				MidiCC int
				CfgKey string
				Type   string
				Value  float32
			}{
				"i0": {
					MidiCC: 30,
					CfgKey: "KICKV",
					Type:   "volume",
				},
			},
		},
		{
			name: "with layers",
			args: args{
				pst: &m.KitPreset{
					Instruments: []m.PresetInstrument{
						{
							Instrument: m.InstrumentRef{
								Controls: map[string]m.Control{
									"pan":   {CfgKey: "RI17P"},
									"pitch": {CfgKey: "RI17T"},
								},
								Layers: map[string]m.Layer{
									"bell": {
										CfgMidiKey: "RI17BKEY",
										Controls: map[string]m.Control{
											"volume": {CfgKey: "RI17BV"},
										},
									},
									"edge": {
										CfgMidiKey: "RI17EKEY",
										Controls: map[string]m.Control{
											"volume": {CfgKey: "RI17EV"},
										},
									},
								},
							},
							Controls: map[string]m.PresetControl{
								"pan":   {MidiCC: 105, Type: "pan"},
								"pitch": {MidiCC: 16, Type: "pitch"},
							},
							Layers: map[string]m.PresetLayer{
								"bell": {
									MidiKey: "ride1_bell",
									Controls: map[string]m.PresetControl{
										"volume": {MidiCC: 104, Type: "volume"},
									},
								},
								"edge": {
									MidiKey: "ride1_edge",
									Controls: map[string]m.PresetControl{
										"volume": {MidiCC: 103, Type: "volume"},
									},
								},
							},
						},
					},
				},
				mididevs: []m.MIDIDevice{
					&MockMMIDIDevice{},
				},
			},
			want: m.KitPreset{
				Instruments: []m.PresetInstrument{
					{
						Instrument: m.InstrumentRef{
							Controls: map[string]m.Control{
								"pan":   {CfgKey: "RI17P"},
								"pitch": {CfgKey: "RI17T"},
							},
							Layers: map[string]m.Layer{
								"bell": {
									CfgMidiKey: "RI17BKEY",
									Controls: map[string]m.Control{
										"volume": {CfgKey: "RI17BV"},
									},
								},
								"edge": {
									CfgMidiKey: "RI17EKEY",
									Controls: map[string]m.Control{
										"volume": {CfgKey: "RI17EV"},
									},
								},
							},
						},
						Controls: map[string]m.PresetControl{
							"pan":   {MidiCC: 105, CfgKey: "RI17P", Type: "pan"},
							"pitch": {MidiCC: 16, CfgKey: "RI17T", Type: "pitch"},
						},
						Layers: map[string]m.PresetLayer{
							"bell": {
								MidiKey:    "ride1_bell",
								CfgMidiKey: "RI17BKEY",
								MidiNote:   53,
								Controls: map[string]m.PresetControl{
									"volume": {MidiCC: 104, CfgKey: "RI17BV", Type: "volume"},
								},
							},
							"edge": {
								MidiKey:    "ride1_edge",
								CfgMidiKey: "RI17EKEY",
								MidiNote:   51,
								Controls: map[string]m.PresetControl{
									"volume": {MidiCC: 103, CfgKey: "RI17EV", Type: "volume"},
								},
							},
						},
					},
				},
			},
			wantErr: false,
			expectedControl: map[string]struct {
				Owner  m.ControlOwner
				MidiCC int
				CfgKey string
				Type   string
				Value  float32
			}{
				"i0": {
					MidiCC: 105,
					CfgKey: "RI17P",
					Type:   "pan",
				},
				"i1": {
					MidiCC: 16,
					CfgKey: "RI17T",
					Type:   "pitch",
				},
				"l2": {
					MidiCC: 104,
					CfgKey: "RI17BV",
					Type:   "volume",
				},
				"l3": {
					MidiCC: 103,
					CfgKey: "RI17EV",
					Type:   "volume",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			preset := tt.args.pst
			if err := preset.PrepareToLoad(tt.args.mididevs); (err != nil) != tt.wantErr {
				t.Errorf("AugmentFromInstrument() error = %v, wantErr %v", err, tt.wantErr)
			}
			if diff := cmp.Diff(tt.want, *preset, cmpopts.IgnoreUnexported(m.KitPreset{}), cmpopts.IgnoreUnexported(m.PresetControl{})); diff != "" {
				t.Errorf("MakeGatewayInfo() mismatch (-want +got):\n%s", diff)
			}

			// Verify controls state using the method from KitPreset
			if diff := m.VerifyControlsForTest(preset, tt.expectedControl); diff != "" {
				t.Error("Controls state does not match expected:", diff)
			}
		})
	}
}

// Real load preset to running linuxsampler
//for listening linuxSampler events:
/* netcat netcat localhost 8888
SUBSCRIBE MISCELLANEOUS
*/
func TestLoadPreset(t *testing.T) {
	ls, err := connectSampler()
	if err != nil {
		t.Errorf("InitSampler() error = %v", err)
		return
	}
	ls.DataDir = path.Join(getProjectPath(), "../_presets")
	osFs := afero.NewOsFs()

	dir := getDBPath()
	d, err := db.NewSqlite(dir)
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	defer d.Close()

	type args struct {
		presetId int64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "kit_preset_1",
			args:    args{presetId: 1},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := LoadPreset(tt.args.presetId, d, ls, osFs)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadPreset() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestLoadPresetToSampler(t *testing.T) {
	type args struct {
		sampler        repo.SamplerRepo
		audDevId       int
		midiDevId      int
		instrumentFile string
	}

	//preconfig
	ls, err := connectSampler()
	if err != nil {
		t.Errorf("InitSampler() error = %v", err)
		return
	}
	aDevId, mDevId, err := InitSampler(ls)
	if err != nil {
		t.Errorf("InitSampler() error = %v", err)
		return
	}

	tests := []struct {
		name     string
		args     args
		wantChnl int
		wantErr  bool
	}{
		{
			"Add channel and load instruments",
			args{
				ls,
				aDevId,
				mDevId,
				"/Users/art/art_work/_projects/artem.brayko/raspidrum/_presets/SMDrums/mappings/ride20/ride20.sfz",
			},
			0,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotChnl, err := LoadPresetToSampler(tt.args.sampler, tt.args.audDevId, tt.args.midiDevId, tt.args.instrumentFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadPresetToSampler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotChnl != tt.wantChnl {
				t.Errorf("LoadPresetToSampler() = %v, want %v", gotChnl, tt.wantChnl)
			}
		})
	}
}
