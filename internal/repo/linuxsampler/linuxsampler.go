package linuxsampler

import (
	"fmt"
	"os"
	"path"

	midi "github.com/raspidrum-srv/internal/app/mididevice"
	m "github.com/raspidrum-srv/internal/model"
	repo "github.com/raspidrum-srv/internal/repo"
	"github.com/raspidrum-srv/internal/repo/file"
	lscp "github.com/raspidrum-srv/libs/liblscp-go"
	"github.com/spf13/afero"
)

var presetDir = "current"
var dirPermission os.FileMode = 0644

// TODO: move to cfg
var sampleRoot = "samples"
var presetRoot = "presets"
var instrumentRoot = "instruments"

// Engine - LinuxSampler engine: gig, sfz, sf2
type LinuxSampler struct {
	Client lscp.Client
	Engine string
}

// Connect
// params grouped by audio channels. Audio channel is key of map
func (l *LinuxSampler) ConnectAudioOutput(driver string, params map[int][]repo.Param[string]) (devId int, err error) {
	devId, err = l.Client.CreateAudioOutputDevice(driver)
	if err != nil {
		return
	}
	if len(params) != 0 {
		// key (k) - channelId
		// value (v) - array of channel params
		for k, v := range params {
			for _, p := range v {
				prm := lscp.Parameter[any]{
					Name:  p.Name,
					Value: p.Value,
				}
				err = l.Client.SetAudioOutputChannelParameter(devId, k, prm)
				if err != nil {
					return
				}
			}
		}
	}
	return
}

// Connect to MIDI port and optional set port parameters (i.e. bindings)
func (l *LinuxSampler) ConnectMidiInput(driver string, params []repo.Param[string]) (devId int, err error) {
	devId, err = l.Client.CreateMidiInputDevice(driver)
	if err != nil {
		return
	}
	if len(params) != 0 {
		for _, p := range params {
			prm := lscp.Parameter[any]{
				Name:  p.Name,
				Value: p.Value,
			}
			err = l.Client.SetMidiInputPortParameter(devId, 0, prm)
			if err != nil {
				return
			}
		}
	}

	return
}

func (l *LinuxSampler) CreateChannel(audioDevId, midiDevId int, instrumentFile string) (channelId int, err error) {
	channelId, err = l.Client.AddSamplerChannel()
	if err != nil {
		return
	}
	err = l.Client.SetChannelAudioOutputDevice(channelId, audioDevId)
	if err != nil {
		return
	}
	err = l.Client.SetChannelMidiInputDevice(channelId, midiDevId)
	if err != nil {
		return
	}
	err = l.Client.LoadSamplerEngine(l.Engine, channelId)
	if err != nil {
		return
	}
	err = l.Client.LoadInstrument(instrumentFile, 0, channelId)
	return
}

func (l *LinuxSampler) LoadPreset(preset *m.KitPreset, mididevs []midi.MIDIDevice, fs afero.Fs) error {

	err := l.genPresetFiles(preset, mididevs, fs)
	if err != nil {
		return fmt.Errorf("failed prepare instrument control files for preset: %w", err)
	}

	// load sfz control files and samples
	// Init "instruments-channel index". Map key - KitPreset.Channel.Key, value - KitPreset.Instruments
	chnlInstr := make(map[string][]int, len(preset.Channels))
	for _, v := range preset.Channels {
		chnlInstr[v.Key] = []int{}
	}
	for i, v := range preset.Instruments {
		// add instrument array index to "instruments-channel index"
		chnlInstr[v.ChannelKey] = append(chnlInstr[v.ChannelKey], i)
	}

	return nil
}

// make sfz control files
func (l *LinuxSampler) genPresetFiles(preset *m.KitPreset, mididevs []midi.MIDIDevice, fs afero.Fs) error {
	presetDir, err := preparePresetDir(presetRoot, fs)
	if err != nil {
		return err
	}
	for _, v := range preset.Instruments {
		// make one file for each instrument
		fcontent := []string{}
		fcontent = append(fcontent, "<control>")
		fcontent = append(fcontent, "default_path="+path.Join(sampleRoot, v.Instrument.Uid, v.Instrument.Key))
		// instrument MIDI Key
		if len(v.MidiKey) > 0 {
			mkeyid, err := mapMidiKey(v.MidiKey, mididevs)
			if err != nil {
				return err
			}
			fcontent = append(fcontent, fmt.Sprintf("#define $%s %d", v.Instrument.CfgMidiKey, mkeyid))
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
				mkeyid, err := mapMidiKey(lv.MidiKey, mididevs)
				if err != nil {
					return err
				}
				fcontent = append(fcontent, fmt.Sprintf("#define $%s %d", lv.CfgMidiKey, mkeyid))
			}
			// layer controls
			for _, lcv := range lv.Controls {
				fcontent = append(fcontent, fmt.Sprintf("#define $%s %d", lcv.CfgKey, lcv.MidiCC))
				fcontent = append(fcontent, fmt.Sprintf("set_cc$%s=%.1f", lcv.CfgKey, lcv.Value))
			}
		}

		intrDir := path.Join(instrumentRoot, v.Instrument.Uid, v.Instrument.Key)
		fcontent = append(fcontent, fmt.Sprintf(`#include "%s.sfz"`, intrDir))
		// save to file
		fname := path.Join(presetDir, v.Instrument.Key+"_ctrl.sfz")
		err = file.WriteLines(fcontent, fname, fs)
		if err != nil {
			return err
		}
	}
	return nil
}

func mapMidiKey(mkey string, mdevs []midi.MIDIDevice) (int, error) {
	devlist := make([]string, len(mdevs))
	for i, d := range mdevs {
		kmap, err := d.GetKeysMapping()
		if err != nil {
			return 0, fmt.Errorf("failed get MIDI Keys mapping for device %s: %w", d.Name(), err)
		}
		devlist[i] = d.Name()
		midiId, ok := kmap[mkey]
		if ok {
			return midiId, nil
		}
	}
	return 0, fmt.Errorf("MIDI devices %s doen't have mapping for MIDI Key %s", devlist, mkey)
}

func preparePresetDir(fpath string, fs afero.Fs) (string, error) {
	dr := path.Join(fpath, presetDir)
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
