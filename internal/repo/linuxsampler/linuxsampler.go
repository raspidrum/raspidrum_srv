package linuxsampler

import (
	"bufio"
	"fmt"
	"os"
	"path"

	midi "github.com/raspidrum-srv/internal/app/mididevice"
	m "github.com/raspidrum-srv/internal/model"
	repo "github.com/raspidrum-srv/internal/repo"
	lscp "github.com/raspidrum-srv/libs/liblscp-go"
	"github.com/spf13/afero"
)

// TODO: move to cfg
var sampleRoot = "samples"
var presetRoot = "presets"
var instrumentRoot = "instruments"
var dirPermission os.FileMode = 0644

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

func (l *LinuxSampler) LoadPreset(preset *m.KitPreset, mididevs []*midi.MIDIDevice, fs afero.Fs) error {

	// Init "instruments-channel index". Map key - KitPreset.Channel.Key, value - KitPreset.Instruments
	chnlInstr := make(map[string][]int, len(preset.Channels))
	for _, v := range preset.Channels {
		chnlInstr[v.Key] = []int{}
	}

	// make sfz control files
	presetDir, err := preparePresetDir(presetRoot, preset.Uid, fs)
	if err != nil {
		return err
	}
	// make file for one instrument
	for i, v := range preset.Instruments {
		// add instrument array index to "instruments-channel index"
		chnlInstr[v.ChannelKey] = append(chnlInstr[v.ChannelKey], i)

		fcontent := []string{}
		fcontent = append(fcontent, "<control>")
		// TODO: add instrument UUID between sampleRoot and key path parts
		fcontent = append(fcontent, "default_path="+path.Join(sampleRoot, v.Instrument.Key))
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
			fcontent = append(fcontent, fmt.Sprintf("set_cc$%s=%f", cv.CfgKey, cv.Value))
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
				fcontent = append(fcontent, fmt.Sprintf("set_cc$%s=%f", lcv.CfgKey, lcv.Value))
			}
		}

		// TODO: add instrument UUID between instrumentRoot and key path parts
		intrDir := path.Join(instrumentRoot, v.Instrument.Key)
		fcontent = append(fcontent, fmt.Sprintf(`#include "%s.sfz"`, intrDir))
		// save to file
		err = writeInstrCtrlFile(fcontent, presetDir, v.Instrument.Key, fs)
		if err != nil {
			return err
		}
	}

	// 5th step: load sfz control files and samples

	return nil
}

func mapMidiKey(mkey string, mdevs []*midi.MIDIDevice) (int, error) {
	devlist := make([]string, len(mdevs))
	for i, d := range mdevs {
		kmap, err := d.GetKeysMapping()
		if err != nil {
			return 0, fmt.Errorf("failed get MIDI Keys mapping for device %s: %w", d.Name, err)
		}
		devlist[i] = d.Name
		midiId, ok := kmap[mkey]
		if ok {
			return midiId, nil
		}
	}
	return 0, fmt.Errorf("MIDI devices %s doen't have mapping for MIDI Key %s", devlist, mkey)
}

func preparePresetDir(fpath string, presetUid string, fs afero.Fs) (string, error) {
	dr := path.Join(fpath, presetUid)
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

func writeInstrCtrlFile(cont []string, fpath string, instrKey string, fs afero.Fs) error {
	fname := path.Join(fpath, instrKey+"_ctrl.sfz")
	f, err := fs.Create(fname)
	if err != nil {
		return fmt.Errorf("failed create file %s: %w", fname, err)
	}
	defer f.Close()

	buf := bufio.NewWriter(f)
	for _, ln := range cont {
		_, err := buf.WriteString(ln + "\n")
		if err != nil {
			return fmt.Errorf("failed write to file %s: %w", fname, err)
		}
	}
	if err := buf.Flush(); err != nil {
		return fmt.Errorf("failed write to file %s: %w", fname, err)
	}
	return nil
}
