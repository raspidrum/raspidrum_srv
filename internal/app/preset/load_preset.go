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
var midiDevices = []midi.MIDIDevice{
	midi.MIDIDevice(&mdev),
}

var osFs = afero.NewOsFs()

// Loads the specified preset into the sampler and returns information about the loaded preset
func LoadPreset(presetId int64, db *d.Sqlite, sampler repo.SamplerRepo, fs afero.Fs) (*m.KitPreset, error) {

	// 1st step: get preset info from db
	pst, err := db.GetPreset(d.ById(presetId))
	if err != nil {
		return nil, fmt.Errorf("failed LoadPreset: %w", err)
	}

	// 2nd step: augment channels and layers info from instrument and instrument preset
	err = augmentFromInstrument(pst)
	if err != nil {
		return nil, err
	}

	// 3rd step: substitute ids of MIDI Keys and MIDI CC
	// skipped: substitute MIDI Keys needed only for generation sfz-ctrl files. MIDI CC stored in db and not needed for substitute

	// 4rd step: load to sampler
	err = sampler.LoadPreset(pst, midiDevices, fs)
	if err != nil {
		return nil, fmt.Errorf("failed load preset to sampler: %w", err)
	}
	return pst, nil
}

// Augment preset controls and layers with data from instrument
func augmentFromInstrument(pst *m.KitPreset) error {
	for i, v := range pst.Instruments {

		// copy instrument control.key to instrument preset
		for kc, vc := range v.Controls {
			ictrl, ok := v.Instrument.Controls[kc]
			if !ok {
				return fmt.Errorf("not found control '%s' in instrument '%s'", kc, v.Instrument.Key)
			}
			vc.CfgKey = ictrl.CfgKey
			pst.Instruments[i].Controls[kc] = vc
		}

		// copy instrument layer MidiKey to instrument preset
		for kl, vl := range v.Layers {
			ilrs, ok := v.Instrument.Layers[kl]
			if !ok {
				return fmt.Errorf("not found layer '%s' in instrument '%s'", kl, v.Instrument.Key)
			}
			vl.CfgMidiKey = ilrs.CfgMidiKey
			// copy instrument layer control.key to instrument preset
			for klc, vlc := range vl.Controls {
				ictrl, ok := ilrs.Controls[klc]
				if !ok {
					return fmt.Errorf("not found control '%s' of layer '%s' in instrument '%s'", klc, kl, v.Instrument.Key)
				}
				vlc.CfgKey = ictrl.CfgKey
				vl.Controls[klc] = vlc
			}
			pst.Instruments[i].Layers[kl] = vl
		}
	}
	return nil
}
