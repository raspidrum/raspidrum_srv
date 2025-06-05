package linuxsampler

import (
	"fmt"
	"io/fs"
	"net"
	"path"
	"reflect"
	"slices"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	m "github.com/raspidrum-srv/internal/model"
	repo "github.com/raspidrum-srv/internal/repo"
	"github.com/raspidrum-srv/internal/repo/file"
	lscp "github.com/raspidrum-srv/libs/liblscp-go"
	"github.com/spf13/afero"
)

func TestLinuxSampler_genPresetFiles(t *testing.T) {
	type fields struct {
		Client lscp.Client
		Engine string
	}
	type args struct {
		preset *m.KitPreset
		fs     afero.Fs
	}
	type res struct {
		dir   string
		files map[string][]string
	}

	rootDir := ""

	tests := []struct {
		name           string
		fields         fields
		args           args
		want           res
		orderImportant bool
		wantErr        bool
	}{
		{
			name: "one instrument w/o layers, w/o controls",
			args: args{
				preset: &m.KitPreset{
					Instruments: []m.PresetInstrument{
						{
							ChannelKey: "1",
							Instrument: m.InstrumentRef{
								Uid:        "1111-ffff",
								Key:        "simple",
								CfgMidiKey: "KEY1",
							},
							MidiKey:  "kick1",
							MidiNote: 36,
						},
					},
					Channels: []m.PresetChannel{
						{Key: "1"},
					},
				},
				fs: afero.NewMemMapFs(),
			},
			orderImportant: true,
			want: res{
				dir: path.Join(rootDir, presetRoot, presetDir),
				files: map[string][]string{
					"simple_ctrl.sfz": {
						"<control>",
						"default_path=samples/1111-ffff/simple/",
						"#define $KEY1 36",
						`#include "instruments/1111-ffff/simple.sfz"`,
					},
					"channel_1.sfz": {
						"#define $VOLMIN 18",
						"#define $VOLSHIFT 24",
						"#define $PITCHMAX 1200",
						"#define $PITCHMIN 600",
						`#include "simple_ctrl.sfz"`,
					},
				},
			},
		},
		{
			name: "two instruments w/o layers, with controls",
			args: args{
				preset: &m.KitPreset{
					Channels: []m.PresetChannel{
						{Key: "1"},
						{Key: "2"},
					},
					Instruments: []m.PresetInstrument{
						{
							ChannelKey: "1",
							Instrument: m.InstrumentRef{
								Uid:        "1111-ffff",
								Key:        "kick",
								CfgMidiKey: "KEYKICK",
							},
							MidiKey:  "kick1",
							MidiNote: 36,
						},
						{
							ChannelKey: "2",
							Instrument: m.InstrumentRef{
								Uid:        "2222-ffff",
								Key:        "snare",
								CfgMidiKey: "KEYSNARE",
							},
							MidiKey:  "snare",
							MidiNote: 38,
							Controls: map[string]*m.PresetControl{
								"volume": {
									CfgKey: "S65NRV",
									MidiCC: 22,
									Value:  95.0,
								},
								"pan": {
									CfgKey: "S65NRP",
									MidiCC: 80,
									Value:  64.0,
								},
							},
						},
					},
				},
				fs: afero.NewMemMapFs(),
			},
			orderImportant: false,
			want: res{
				dir: path.Join(rootDir, presetRoot, presetDir),
				files: map[string][]string{
					"kick_ctrl.sfz": {
						"<control>",
						"default_path=samples/1111-ffff/kick/",
						"#define $KEYKICK 36",
						`#include "instruments/1111-ffff/kick.sfz"`,
					},
					"snare_ctrl.sfz": {
						"<control>",
						"default_path=samples/2222-ffff/snare/",
						"#define $KEYSNARE 38",
						"#define $S65NRV 22",
						"set_cc$S65NRV=95.0",
						"#define $S65NRP 80",
						"set_cc$S65NRP=64.0",
						`#include "instruments/2222-ffff/snare.sfz"`,
					},
					"channel_1.sfz": {
						"#define $VOLMIN 18",
						"#define $VOLSHIFT 24",
						"#define $PITCHMAX 1200",
						"#define $PITCHMIN 600",
						`#include "kick_ctrl.sfz"`,
					},
					"channel_2.sfz": {
						"#define $VOLMIN 18",
						"#define $VOLSHIFT 24",
						"#define $PITCHMAX 1200",
						"#define $PITCHMIN 600",
						`#include "snare_ctrl.sfz"`,
					},
				},
			},
		},
		{
			name: "two instruments with layers, with controls",
			args: args{
				preset: &m.KitPreset{
					Channels: []m.PresetChannel{
						{Key: "1"},
					},
					Instruments: []m.PresetInstrument{
						{
							ChannelKey: "1",
							Instrument: m.InstrumentRef{
								Uid: "1111-ffff",
								Key: "ride",
							},
							Layers: map[string]m.PresetLayer{
								"bell": {
									CfgMidiKey: "RI17BKEY",
									MidiKey:    "ride1_bell",
									MidiNote:   53,
									Controls: map[string]*m.PresetControl{
										"volume": {
											CfgKey: "RI17BV",
											MidiCC: 104,
											Value:  95.0,
										},
									},
								},
								"edge": {
									CfgMidiKey: "RI17EKEY",
									MidiKey:    "ride1_edge",
									MidiNote:   51,
									Controls: map[string]*m.PresetControl{
										"volume": {
											CfgKey: "RI17EV",
											MidiCC: 103,
											Value:  95.0,
										},
									},
								},
							},
						},
						{
							ChannelKey: "1",
							Instrument: m.InstrumentRef{
								Uid:        "2222-ffff",
								Key:        "snare",
								CfgMidiKey: "KEYSNARE",
							},
							MidiKey:  "snare",
							MidiNote: 38,
							Controls: map[string]*m.PresetControl{
								"volume": {
									CfgKey: "S65NRV",
									MidiCC: 22,
									Value:  95.0,
								},
								"pan": {
									CfgKey: "S65NRP",
									MidiCC: 80,
									Value:  64.0,
								},
							},
						},
					},
				},
				fs: afero.NewMemMapFs(),
			},
			orderImportant: false,
			want: res{
				dir: path.Join(rootDir, presetRoot, presetDir),
				files: map[string][]string{
					"ride_ctrl.sfz": {
						"<control>",
						"default_path=samples/1111-ffff/ride/",
						"#define $RI17BKEY 53",
						"#define $RI17BV 104",
						"set_cc$RI17BV=95.0",
						"#define $RI17EKEY 51",
						"#define $RI17EV 103",
						"set_cc$RI17EV=95.0",
						`#include "instruments/1111-ffff/ride.sfz"`,
					},
					"snare_ctrl.sfz": {
						"<control>",
						"default_path=samples/2222-ffff/snare/",
						"#define $KEYSNARE 38",
						"#define $S65NRV 22",
						"set_cc$S65NRV=95.0",
						"#define $S65NRP 80",
						"set_cc$S65NRP=64.0",
						`#include "instruments/2222-ffff/snare.sfz"`,
					},
					"channel_1.sfz": {
						"#define $VOLMIN 18",
						"#define $VOLSHIFT 24",
						"#define $PITCHMAX 1200",
						"#define $PITCHMIN 600",
						`#include "ride_ctrl.sfz"`,
						`#include "snare_ctrl.sfz"`,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &LinuxSampler{
				Client:  tt.fields.Client,
				Engine:  tt.fields.Engine,
				DataDir: rootDir,
			}
			if _, err := l.genPresetFiles(tt.args.preset, tt.args.fs); (err != nil) != tt.wantErr {
				t.Errorf("LinuxSampler.LoadPreset() error = %v, wantErr %v", err, tt.wantErr)
			}
			// get files from MemFs
			cont, err := readResultFiles(tt.want.dir, tt.args.fs)
			if err != nil {
				t.Errorf("LinuxSampler.LoadPreset() error = failed get generated ctrl files: %v", err)
			}
			// compare result with want
			errs := map[string][]string{}
			for k, v := range tt.want.files {
				gotCont, ok := cont[k]
				if !ok {
					errs[k] = []string{fmt.Sprintf("absent file %s", k)}
					continue
				}
				var diff []string
				if tt.orderImportant {
					diff = compareLines(v, gotCont)
				} else {
					slices.Sort(v)
					slices.Sort(gotCont)
					diff = compareLines(v, gotCont)
				}
				if len(diff) > 0 {
					errs[k] = diff
				}
			}
			//find files in result that absent in want.files
			for k := range cont {
				_, ok := tt.want.files[k]
				if ok { // file has been processed in loop above
					continue
				}
				errs[k] = []string{fmt.Sprintf(`got file "%s" but not wanted`, k)}
			}

			if len(errs) > 0 {
				t.Errorf(`LinuxSampler.LoadPreset() error = files content diff: %v`, errs)
			}
		})
	}
}

func readResultFiles(dname string, afs afero.Fs) (map[string][]string, error) {
	if ok, err := afero.DirExists(afs, dname); !ok {
		return nil, fmt.Errorf(`dir "%s" absent: %w`, dname, err)
	}
	res := make(map[string][]string)
	err := afero.Walk(afs, dname, func(filepath string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		fcont, err := file.ReadLines(filepath, afs)
		if err != nil {
			return err
		}
		_, fname := path.Split(filepath)
		res[fname] = fcont
		return nil
	})
	return res, err
}

func compareLines(want []string, got []string) []string {
	res := []string{}
	a_len := len(want)
	b_len := len(got)
	for i, v := range want {
		if i > b_len-1 {
			res = append(res, fmt.Sprintf(`want: "%s"\n got:""`, v))
			continue
		}
		if v != got[i] {
			res = append(res, fmt.Sprintf(`want: "%s"\n got:"%s"`, v, got[i]))
		}
	}
	if a_len < b_len {
		for i := a_len; i < a_len+b_len-a_len; i++ {
			res = append(res, fmt.Sprintf(`want: ""\n got:"%s"`, got[i]))
		}
	}
	return res
}

func Test_compareLines(t *testing.T) {
	type args struct {
		a []string
		b []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "fully equal",
			args: args{
				a: []string{
					"a1",
					"a2",
				},
				b: []string{
					"a1",
					"a2",
				},
			},
			want: []string{},
		},
		{
			name: "lens equal, content diff",
			args: args{
				a: []string{
					"a1",
					"a2",
				},
				b: []string{
					"b1",
					"b2",
				},
			},
			want: []string{
				`want: "a1"\n got:"b1"`,
				`want: "a2"\n got:"b2"`,
			},
		},
		{
			name: "a > b, content diff",
			args: args{
				a: []string{
					"a1",
					"a2",
					"a3",
				},
				b: []string{
					"b1",
					"b2",
				},
			},
			want: []string{
				`want: "a1"\n got:"b1"`,
				`want: "a2"\n got:"b2"`,
				`want: "a3"\n got:""`,
			},
		},
		{
			name: "a < b, content diff",
			args: args{
				a: []string{
					"a1",
					"a2",
				},
				b: []string{
					"b1",
					"b2",
					"b3",
				},
			},
			want: []string{
				`want: "a1"\n got:"b1"`,
				`want: "a2"\n got:"b2"`,
				`want: ""\n got:"b3"`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := compareLines(tt.args.a, tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("compareLines() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLinuxSampler_loadToSampler(t *testing.T) {
	rootDir := ""
	type fields struct {
		Client lscp.Client
		Engine string
	}
	type args struct {
		preset          *m.KitPreset
		instrumentFiles map[string]string
	}
	type res struct {
		lscpCommands []string
		channels     repo.SamplerChannels
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    res
		wantErr bool
	}{
		{
			name: "one channel, one instrument",
			fields: fields{
				Client: lscp.NewClient("pipe", "0", "1s"),
				Engine: "sfz",
			},
			args: args{
				preset: &m.KitPreset{
					Channels: []m.PresetChannel{
						{Key: "1",
							Controls: map[string]*m.PresetControl{
								"volume": {
									Type:  "volume",
									Value: 1.0,
								},
							},
						},
					},
					Instruments: []m.PresetInstrument{
						{
							Instrument: m.InstrumentRef{
								Uid: "1111-ffff",
								Key: "simple",
							},
							ChannelKey: "1",
						},
					},
				},
				instrumentFiles: map[string]string{
					"1111-ffff": path.Join(rootDir, presetRoot, presetDir, "simple_ctrl.sfz"),
					"channel_1": path.Join(rootDir, presetRoot, presetDir, "channel_1.sfz"),
				},
			},
			want: res{
				lscpCommands: []string{
					"ADD CHANNEL",
					"SET CHANNEL AUDIO_OUTPUT_DEVICE 0 0",
					"SET CHANNEL MIDI_INPUT_DEVICE 0 0",
					"LOAD ENGINE sfz 0",
					"LOAD INSTRUMENT '" + path.Join(rootDir, presetRoot, presetDir, "channel_1.sfz") + "' 0 0",
					"SET CHANNEL VOLUME 0 1.00",
				},
				channels: repo.SamplerChannels{
					"1": 0,
				},
			},
			wantErr: false,
		},
		{
			name: "one channel, two instrument",
			fields: fields{
				Client: lscp.NewClient("pipe", "0", "1s"),
				Engine: "sfz",
			},
			args: args{
				preset: &m.KitPreset{
					Channels: []m.PresetChannel{
						{Key: "1"},
					},
					Instruments: []m.PresetInstrument{
						{
							Instrument: m.InstrumentRef{
								Uid: "1111-ffff",
								Key: "kick",
							},
							ChannelKey: "1",
						},
						{
							Instrument: m.InstrumentRef{
								Uid: "2222-ffff",
								Key: "snare",
							},
							ChannelKey: "1",
						},
					},
				},
				instrumentFiles: map[string]string{
					"1111-ffff": path.Join(rootDir, presetRoot, presetDir, "kick_ctrl.sfz"),
					"2222-ffff": path.Join(rootDir, presetRoot, presetDir, "snare_ctrl.sfz"),
					"channel_1": path.Join(rootDir, presetRoot, presetDir, "channel_1.sfz"),
				},
			},
			want: res{
				lscpCommands: []string{
					"ADD CHANNEL",
					"SET CHANNEL AUDIO_OUTPUT_DEVICE 0 0",
					"SET CHANNEL MIDI_INPUT_DEVICE 0 0",
					"LOAD ENGINE sfz 0",
					"LOAD INSTRUMENT '" + path.Join(rootDir, presetRoot, presetDir, "channel_1.sfz") + "' 0 0",
				},
				channels: repo.SamplerChannels{
					"1": 0,
				},
			},
			wantErr: false,
		},
		{
			name: "two channel, two instrument",
			fields: fields{
				Client: lscp.NewClient("pipe", "0", "1s"),
				Engine: "sfz",
			},
			args: args{
				preset: &m.KitPreset{
					Channels: []m.PresetChannel{
						{Key: "1"},
						{Key: "2"},
					},
					Instruments: []m.PresetInstrument{
						{
							Instrument: m.InstrumentRef{
								Uid: "1111-ffff",
								Key: "kick",
							},
							ChannelKey: "2",
						},
						{
							Instrument: m.InstrumentRef{
								Uid: "2222-ffff",
								Key: "snare",
							},
							ChannelKey: "1",
						},
					},
				},
				instrumentFiles: map[string]string{
					"1111-ffff": path.Join(rootDir, presetRoot, presetDir, "kick_ctrl.sfz"),
					"2222-ffff": path.Join(rootDir, presetRoot, presetDir, "snare_ctrl.sfz"),
					"channel_1": path.Join(rootDir, presetRoot, presetDir, "channel_1.sfz"),
					"channel_2": path.Join(rootDir, presetRoot, presetDir, "channel_2.sfz"),
				},
			},
			want: res{
				lscpCommands: []string{
					"ADD CHANNEL",
					"SET CHANNEL AUDIO_OUTPUT_DEVICE 0 0",
					"SET CHANNEL MIDI_INPUT_DEVICE 0 0",
					"LOAD ENGINE sfz 0",
					"LOAD INSTRUMENT '" + path.Join(rootDir, presetRoot, presetDir, "channel_1.sfz") + "' 0 0",
					"ADD CHANNEL",
					"SET CHANNEL AUDIO_OUTPUT_DEVICE 1 0",
					"SET CHANNEL MIDI_INPUT_DEVICE 1 0",
					"LOAD ENGINE sfz 1",
					"LOAD INSTRUMENT '" + path.Join(rootDir, presetRoot, presetDir, "channel_2.sfz") + "' 0 1",
				},
				channels: repo.SamplerChannels{
					"1": 0,
					"2": 1,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// start mock server
			clientConn, serverConn := net.Pipe()
			tt.fields.Client.Conn = clientConn
			mockServer := startMockPipeServer(serverConn)
			l := &LinuxSampler{
				Client:  tt.fields.Client,
				Engine:  tt.fields.Engine,
				DataDir: rootDir,
			}

			gotChannels, err := l.loadToSampler(0, 0, tt.args.preset, tt.args.instrumentFiles)
			if (err != nil) != tt.wantErr {
				t.Errorf("LinuxSampler.loadToSampler() error = %v, wantErr %v", err, tt.wantErr)
			}
			// Let some time for finish processing and then stop
			time.Sleep(100 * time.Millisecond)
			mockServer.stop()
			// compare lscp commends
			if !reflect.DeepEqual(mockServer.getMessages(), tt.want.lscpCommands) {
				t.Errorf("diff lscp commands: %v", compareLines(tt.want.lscpCommands, mockServer.getMessages()))
			}
			// compare channels
			if diff := cmp.Diff(tt.want.channels, gotChannels); diff != "" {
				t.Errorf("mismatch channels (-want +got):\n%s", diff)
			}

			clientConn.Close()
			serverConn.Close()
		})
	}
}
