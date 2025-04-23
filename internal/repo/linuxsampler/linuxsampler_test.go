package linuxsampler

import (
	"fmt"
	"io/fs"
	"path"
	"reflect"
	"testing"

	midi "github.com/raspidrum-srv/internal/app/mididevice"
	m "github.com/raspidrum-srv/internal/model"
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
		preset   *m.KitPreset
		mididevs []*midi.MIDIDevice
		fs       afero.Fs
	}
	type res struct {
		dir   string
		files map[string][]string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    res
		wantErr bool
	}{
		{
			name: "one instrument w/o layers, w/o controls in one channel",
			args: args{
				preset: &m.KitPreset{
					Uid: "aaaa-bbbb-cccc-dddd",
					Channels: []m.PresetChannel{
						{Key: "1"},
					},
					Instruments: []m.PresetInstrument{
						{
							Instrument: m.InstrumentRef{
								Uid:        "1111-ffff",
								Key:        "simple",
								CfgMidiKey: "KEY1",
							},
							ChannelKey: "1",
							MidiKey:    "kick",
						},
					},
				},
				mididevs: []*midi.MIDIDevice{
					// TODO: mock GetKeysMapping
					{
						DevId: "0:0",
						Name:  "Dummy",
					},
				},
				fs: afero.NewMemMapFs(),
			},
			want: res{
				dir: path.Join(presetRoot, presetDir),
				files: map[string][]string{
					"simple_ctrl.sfz": {
						"<control>",
						"default_path=samples/1111-ffff/simple",
						"#define $KEY1 36",
						`#include "instruments/1111-ffff/simple.sfz"`,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &LinuxSampler{
				Client: tt.fields.Client,
				Engine: tt.fields.Engine,
			}
			if err := l.genPresetFiles(tt.args.preset, tt.args.mididevs, tt.args.fs); (err != nil) != tt.wantErr {
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
				diff := compareLines(v, gotCont)
				if len(diff) > 0 {
					errs[k] = diff
				}
			}
			// TODO: loop for files in result and check existance in want.files

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
