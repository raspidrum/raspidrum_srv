package linuxsampler

import (
	repo "github.com/raspidrum-srv/internal/repo"
	lscp "github.com/raspidrum-srv/libs/liblscp-go"
)

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
