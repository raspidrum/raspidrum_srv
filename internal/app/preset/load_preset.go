package preset

import (
	"fmt"

	"github.com/spf13/afero"

	midi "github.com/raspidrum-srv/internal/app/mididevice"
	m "github.com/raspidrum-srv/internal/model"
	"github.com/raspidrum-srv/internal/repo"
	d "github.com/raspidrum-srv/internal/repo/db"
)

// TODO: init MIDI device on connect/reconnect (and startup)
var mdev = midi.NewUSBMIDIDevice("0:0", "Dummy")
var midiDevices = []m.MIDIDevice{
	&mdev,
}

// Loads the specified preset into the sampler and returns information about the loaded preset
func LoadPreset(presetId int64, db *d.Sqlite, sampler repo.SamplerRepo, fs afero.Fs) (*m.KitPreset, error) {

	// 1st step: get preset info from db
	pst, err := db.GetPreset(d.ById(presetId))
	if err != nil {
		return nil, fmt.Errorf("failed LoadPreset: %w", err)
	}

	// 2nd step: augment channels and layers info from instrument and instrument preset
	err = pst.AugmentAndIndex(midiDevices)
	if err != nil {
		return nil, err
	}

	// 3rd step: substitute ids of MIDI Keys and MIDI CC
	// skipped: substitute MIDI Keys needed only for generation sfz-ctrl files. MIDI CC stored in db and not needed for substitute

	// 4rd step: init sampler
	audioDevId, midiDevId, err := InitSampler(sampler)
	if err != nil {
		return nil, fmt.Errorf("failed init sampler: %w", err)
	}

	// 5th step: load to sampler
	err = sampler.LoadPreset(audioDevId, midiDevId, pst, fs)
	if err != nil {
		return nil, fmt.Errorf("failed load preset to sampler: %w", err)
	}
	return pst, nil
}

// deprecated
// old func for testing load one instrument-file in new sampler channel
func LoadPresetToSampler(sampler repo.SamplerRepo, audDevId, midiDevId int, instrumentFile string) (chnl int, err error) {
	chnl, err = sampler.CreateChannel(audDevId, midiDevId)
	if err != nil {
		return chnl, fmt.Errorf("failed load preset: %w", err)
	}
	return chnl, sampler.LoadInstrument(instrumentFile, 0, chnl)
}
