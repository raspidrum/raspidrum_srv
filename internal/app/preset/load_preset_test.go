package preset

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	midi "github.com/raspidrum-srv/internal/app/mididevice"
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
	}, nil
}

func Test_augmentFromInstrument(t *testing.T) {
	type args struct {
		pst      *m.KitPreset
		mididevs []midi.MIDIDevice
	}
	tests := []struct {
		name    string
		args    args
		want    m.KitPreset
		wantErr bool
	}{
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
								},
							},
						},
					},
				},
				mididevs: []midi.MIDIDevice{
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
							},
						},
					},
				},
			},
			wantErr: false,
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
								"pan":   {MidiCC: 105},
								"pitch": {MidiCC: 16},
							},
							Layers: map[string]m.PresetLayer{
								"bell": {
									MidiKey: "ride1_bell",
									Controls: map[string]m.PresetControl{
										"volume": {MidiCC: 104},
									},
								},
								"edge": {
									MidiKey: "ride1_edge",
									Controls: map[string]m.PresetControl{
										"volume": {MidiCC: 103},
									},
								},
							},
						},
					},
				},
				mididevs: []midi.MIDIDevice{
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
							"pan":   {MidiCC: 105, CfgKey: "RI17P"},
							"pitch": {MidiCC: 16, CfgKey: "RI17T"},
						},
						Layers: map[string]m.PresetLayer{
							"bell": {
								MidiKey:    "ride1_bell",
								CfgMidiKey: "RI17BKEY",
								MidiNote:   53,
								Controls: map[string]m.PresetControl{
									"volume": {MidiCC: 104, CfgKey: "RI17BV"},
								},
							},
							"edge": {
								MidiKey:    "ride1_edge",
								CfgMidiKey: "RI17EKEY",
								MidiNote:   51,
								Controls: map[string]m.PresetControl{
									"volume": {MidiCC: 103, CfgKey: "RI17EV"},
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			preset := tt.args.pst
			if err := augmentFromInstrument(preset, tt.args.mididevs); (err != nil) != tt.wantErr {
				t.Errorf("augmentFromInstrument() error = %v, wantErr %v", err, tt.wantErr)
			}
			if diff := cmp.Diff(tt.want, *preset); diff != "" {
				t.Errorf("MakeGatewayInfo() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
