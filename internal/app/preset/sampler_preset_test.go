package preset

import (
	"fmt"
	"testing"

	"github.com/raspidrum-srv/internal/repo"
	lsampler "github.com/raspidrum-srv/internal/repo/linuxsampler"
	lscp "github.com/raspidrum-srv/libs/liblscp-go"
)

func connectSampler() (*lsampler.LinuxSampler, error) {
	lsClient := lscp.NewClient("localhost", "8888", "1s")

	err := lsClient.Connect()
	ls := lsampler.LinuxSampler{
		Client: lsClient,
		Engine: "sfz",
	}
	if err != nil {
		return &ls, fmt.Errorf("Failed connect to LinuxSampler: %v", err)
	}
	return &ls, nil
}

func TestInitSampler(t *testing.T) {
	type args struct {
		sampler repo.SamplerRepo
	}

	// preconfig
	ls, err := connectSampler()
	if err != nil {
		t.Errorf("InitSampler() error = %v", err)
		return
	}

	tests := []struct {
		name           string
		sampler        repo.SamplerRepo
		wantAudioDevId int
		wantMidiDevId  int
		wantErr        bool
	}{
		{
			"init sampler",
			ls,
			0,
			0,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAudioDevId, gotMidiDevId, err := InitSampler(tt.sampler)
			if (err != nil) != tt.wantErr {
				t.Errorf("InitSampler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotAudioDevId != tt.wantAudioDevId {
				t.Errorf("InitSampler() gotAudioDevId = %v, want %v", gotAudioDevId, tt.wantAudioDevId)
			}
			if gotMidiDevId != tt.wantMidiDevId {
				t.Errorf("InitSampler() gotMidiDevId = %v, want %v", gotMidiDevId, tt.wantMidiDevId)
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
