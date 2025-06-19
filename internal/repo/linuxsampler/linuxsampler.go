package linuxsampler

import (
	repo "github.com/raspidrum-srv/internal/repo"
	lscp "github.com/raspidrum-srv/libs/liblscp-go"
)

// Engine - LinuxSampler engine: gig, sfz, sf2
type LinuxSampler struct {
	Client  lscp.Client
	Engine  string
	DataDir string // root dir for sfz-files, samples and presets
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

func (l *LinuxSampler) CreateChannel(audioDevId, midiDevId int) (channelId int, err error) {
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
	return
}

func (l *LinuxSampler) LoadInstrument(instrumentFile string, instrIdx int, channelId int) error {
	return l.Client.LoadInstrument(instrumentFile, 0, channelId)
}

func (l *LinuxSampler) SetChannelVolume(samplerChn int, volume float32) error {
	return l.Client.SetChannelVolume(samplerChn, volume)
}

func (l *LinuxSampler) SendMidiCC(samplerChn int, cc int, value float32) error {
	return l.Client.SendChannelMidiData(samplerChn, "CC", cc, int(value))
}

func (l *LinuxSampler) SetGlobalVolume(volume float32) error {
	return l.Client.SetVolume(volume)
}
