//go:build darwin
// +build darwin

package coremidi

import (
	"github.com/raspidrum-srv/internal/app/midi"
)

type CoreMidiProvider struct {
}

func NewCoreMidiProvider() *CoreMidiProvider {
	return &CoreMidiProvider{}
}

func (p *CoreMidiProvider) GetMIDIPorts() ([]midi.MIDIPortInfo, error) {
	return []midi.MIDIPortInfo{
		{
			Driver:     "COREMIDI",
			DevId:      "0:0",
			PortId:     "vmpk vmpk out",
			Name:       "CoreMIDI Device",
			State:      midi.MIDIPortStateConnected,
			DeviceType: midi.MIDIPortTypeSoftware,
		},
	}, nil
}

func (p *CoreMidiProvider) SubscribeDeviceState(listener func(devId string, state midi.MIDIPortState)) error {
	return nil
}

func (p *CoreMidiProvider) UnsubscribeDeviceState() error {
	return nil
}
