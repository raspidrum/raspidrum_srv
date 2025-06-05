package linuxsampler

import (
	"fmt"
	"os"
	"path"

	m "github.com/raspidrum-srv/internal/model"
	repo "github.com/raspidrum-srv/internal/repo"
	"github.com/raspidrum-srv/internal/repo/file"
	"github.com/spf13/afero"
)

var presetDir = "current"
var dirPermission os.FileMode = os.ModePerm

// TODO: move to cfg
var sampleRoot = "samples"
var presetRoot = "presets"
var instrumentRoot = "instruments"

func (l *LinuxSampler) LoadPreset(audioDevId, midiDevId int, preset *m.KitPreset, fs afero.Fs) (repo.SamplerChannels, error) {

	instrFiles, err := l.genPresetFiles(preset, fs)
	if err != nil {
		return nil, fmt.Errorf("failed prepare instrument control files for preset: %w", err)
	}

	// load sfz control files and samples in sampler
	chnls, err := l.loadToSampler(audioDevId, midiDevId, preset, instrFiles)
	if err != nil {
		return nil, fmt.Errorf("failed load instrument config and samples to sampler: %w", err)
	}
	return chnls, nil
}

// make sfz control files
// return map instrument.uid : filename with path
func (l *LinuxSampler) genPresetFiles(preset *m.KitPreset, fs afero.Fs) (map[string]string, error) {
	presetFiles := map[string]string{}
	presetDir, err := preparePresetDir(l.DataDir, fs)
	if err != nil {
		return nil, err
	}
	sampleDir := path.Join(l.DataDir, sampleRoot)

	//content for channel files - used to load many instrument files in one sampler channel
	chnlFiles := make(map[string][]string, len(preset.Channels))
	for _, v := range preset.Channels {
		chnlFiles[v.Key] = []string{}
	}

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
		instrName := v.Instrument.Key + "_ctrl.sfz"
		fname := path.Join(presetDir, instrName)
		presetFiles[v.Instrument.Uid] = fname
		err = file.WriteLines(fcontent, fname, fs)
		if err != nil {
			return nil, err
		}

		chnlFiles[v.ChannelKey] = append(chnlFiles[v.ChannelKey], fmt.Sprintf(`#include "%s"`, instrName))
	}

	// write channel files
	for k, v := range chnlFiles {
		chnlName := "channel_" + k
		fname := path.Join(presetDir, fmt.Sprintf("%s.sfz", chnlName))
		cont := []string{}
		cont = append(cont, getControlLimits()...)
		cont = append(cont, v...)
		err = file.WriteLines(cont, fname, fs)
		if err != nil {
			return nil, err
		}
		presetFiles[chnlName] = fname
	}

	return presetFiles, nil
}

// make define sfzvariables for control limits
func getControlLimits() []string {
	return []string{
		"#define $VOLMIN 18",
		"#define $VOLSHIFT 24",
		"#define $PITCHMAX 1200",
		"#define $PITCHMIN 600",
	}
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
func (l *LinuxSampler) loadToSampler(audDevId, midiDevId int, preset *m.KitPreset, instrfiles map[string]string) (repo.SamplerChannels, error) {
	channels := repo.SamplerChannels{}

	//loading instruments
	for _, cv := range preset.Channels {
		// create channel
		chnlId, err := l.CreateChannel(audDevId, midiDevId)
		if err != nil {
			return nil, fmt.Errorf("failed create sampler channel: %w", err)
		}
		channels[cv.Key] = chnlId

		// load instruments to channel
		chnlName := "channel_" + cv.Key
		fname, ok := instrfiles[chnlName]
		if !ok {
			return nil, fmt.Errorf("failed load instrument: not found filename for channel instruments file %s", chnlName)
		}
		err = l.LoadInstrument(fname, 0, chnlId)
		if err != nil {
			return nil, fmt.Errorf("failed load instruments %s to sampler: %w", chnlName, err)
		}

		// set channel controls
		for _, ccv := range cv.Controls {
			if len(ccv.CfgKey) == 0 && m.ControlTypeFromString[ccv.Type] == m.CTVolume {
				l.SetChannelVolume(chnlId, ccv.Value)
			}
		}
	}

	return channels, nil
}
