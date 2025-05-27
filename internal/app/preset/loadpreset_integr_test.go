//go:build integration

package preset

import (
	"path"
	"testing"

	"github.com/raspidrum-srv/internal/repo"
	"github.com/raspidrum-srv/internal/repo/db"
	"github.com/spf13/afero"
)

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
