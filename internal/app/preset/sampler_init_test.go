//go:build integration

package preset

import (
	"testing"

	"github.com/raspidrum-srv/internal/repo"
)

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
