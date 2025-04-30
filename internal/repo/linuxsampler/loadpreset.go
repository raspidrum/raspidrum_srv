package linuxsampler

import (
	"fmt"
	"os"
	"path"

	m "github.com/raspidrum-srv/internal/model"
	"github.com/raspidrum-srv/internal/repo/file"
	"github.com/spf13/afero"
)

var presetDir = "current"
var dirPermission os.FileMode = os.ModePerm

// TODO: move to cfg
var sampleRoot = "samples"
var presetRoot = "presets"
var instrumentRoot = "instruments"

func (l *LinuxSampler) LoadPreset(audioDevId, midiDevId int, preset *m.KitPreset, fs afero.Fs) error {

	instrFiles, err := l.genPresetFiles(preset, fs)
	if err != nil {
		return fmt.Errorf("failed prepare instrument control files for preset: %w", err)
	}

	// load sfz control files and samples in sampler
	_, err = l.loadToSampler(audioDevId, midiDevId, preset, instrFiles)
	if err != nil {
		return fmt.Errorf("failed load instrument config and samples to sampler: %w", err)
	}
	return nil
}

// make sfz control files
// return map instrument.uid : filename with path
func (l *LinuxSampler) genPresetFiles(preset *m.KitPreset, fs afero.Fs) (map[string]string, error) {
	instrFiles := map[string]string{}
	presetDir, err := preparePresetDir(l.DataDir, fs)
	if err != nil {
		return nil, err
	}
	sampleDir := path.Join(l.DataDir, sampleRoot)
	for _, v := range preset.Instruments {
		// make one file for each instrument
		fcontent := []string{}
		fcontent = append(fcontent, "<control>")
		fcontent = append(fcontent, "default_path="+path.Join(sampleDir, v.Instrument.Uid, v.Instrument.Key)+"/")
		// instrument MIDI Key
		if len(v.MidiKey) > 0 {
			fcontent = append(fcontent, fmt.Sprintf("#define $%s %d", v.Instrument.CfgMidiKey, v.MidiNote))
		}
		// instrument Controls
		for _, cv := range v.Controls {
			fcontent = append(fcontent, fmt.Sprintf("#define $%s %d", cv.CfgKey, cv.MidiCC))
			fcontent = append(fcontent, fmt.Sprintf("set_cc$%s=%.1f", cv.CfgKey, cv.Value))
		}

		// instrument layers
		for _, lv := range v.Layers {
			// layer MIDI Key
			if len(lv.MidiKey) > 0 {
				fcontent = append(fcontent, fmt.Sprintf("#define $%s %d", lv.CfgMidiKey, lv.MidiNote))
			}
			// layer controls
			for _, lcv := range lv.Controls {
				fcontent = append(fcontent, fmt.Sprintf("#define $%s %d", lcv.CfgKey, lcv.MidiCC))
				fcontent = append(fcontent, fmt.Sprintf("set_cc$%s=%.1f", lcv.CfgKey, lcv.Value))
			}
		}

		intrDir := path.Join(l.DataDir, instrumentRoot, v.Instrument.Uid, v.Instrument.Key)
		fcontent = append(fcontent, fmt.Sprintf(`#include "%s.sfz"`, intrDir))
		// save to file
		fname := path.Join(presetDir, v.Instrument.Key+"_ctrl.sfz")
		instrFiles[v.Instrument.Uid] = fname
		err = file.WriteLines(fcontent, fname, fs)
		if err != nil {
			return nil, err
		}
	}
	return instrFiles, nil
}

// recreate dir with instrument control files (<instrument>_ctrl.sfz)
func preparePresetDir(rootDir string, fs afero.Fs) (string, error) {
	dr := path.Join(rootDir, presetRoot, presetDir)
	err := fs.RemoveAll(dr)
	if err != nil {
		return dr, fmt.Errorf("failed to prepare preset directory %s: %w", dr, err)
	}
	err = fs.MkdirAll(dr, dirPermission)
	if err != nil {
		return dr, fmt.Errorf("failed to prepare preset directory %s: %w", dr, err)
	}
	return dr, nil
}

// Create sampler channels and load into its preset instruments
// return map: key - Channel.Key, value - sampler channel Id
func (l *LinuxSampler) loadToSampler(audDevId, midiDevId int, preset *m.KitPreset, instrfiles map[string]string) (map[string]int, error) {
	channels := map[string]int{}

	// Init "instruments-channel index". Map key - KitPreset.Channel.Key, value - KitPreset.Instruments index
	chnlInstr := make(map[string][]int, len(preset.Channels))
	for _, v := range preset.Channels {
		chnlInstr[v.Key] = []int{}
	}
	for i, v := range preset.Instruments {
		// add instrument array index to "instruments-channel index"
		chnlInstr[v.ChannelKey] = append(chnlInstr[v.ChannelKey], i)
	}

	//loading instruments
	for _, cv := range preset.Channels {
		// create channel
		chnlId, err := l.CreateChannel(audDevId, midiDevId)
		if err != nil {
			return nil, fmt.Errorf("failed create sampler channel: %w", err)
		}
		channels[cv.Key] = chnlId
		cins := chnlInstr[cv.Key]

		// load instruments to channel
		for _, iidx := range cins {
			instr := preset.Instruments[iidx]
			fname, ok := instrfiles[instr.Instrument.Uid]
			if !ok {
				return nil, fmt.Errorf("failed load instrument: not found filename for instrument %s", instr.Instrument.Key)
			}
			err = l.LoadInstrument(fname, 0, chnlId)
			if err != nil {
				return nil, fmt.Errorf("failed load instrument %s to sampler: %w", instr.Instrument.Key, err)
			}
		}

		// set channel controls
		for _, ccv := range cv.Controls {
			// TODO: extract control types to enum
			if len(ccv.CfgKey) == 0 && ccv.Type == "volume" {
				l.SetChannelVolume(chnlId, float64(ccv.Value))
			}
		}
	}

	return channels, nil
}
