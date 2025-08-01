package midi

import (
	"fmt"
	"log/slog"
)

type MIDIDeviceProvider interface {
	GetMIDIPorts() ([]MIDIPortInfo, error)
	SubscribeDeviceState(listener func(devId string, state MIDIPortState)) error
	UnsubscribeDeviceState() error
}

type MIDIDevice interface {
	GetKeysMapping() (map[string]int, error)
	GetOutPorts(isConnected bool) ([]MIDIPortInfo, error)
	Name() string
}

type usbMIDIDevice struct {
	provider MIDIDeviceProvider
	outPorts MIDIPorts
}

func NewMIDIDevice(provider MIDIDeviceProvider) (MIDIDevice, error) {
	if provider == nil {
		return nil, fmt.Errorf("provider cannot be nil")
	}
	m := &usbMIDIDevice{
		provider: provider,
	}
	provider.SubscribeDeviceState(m.providerEventHandler)
	return m, nil
}

func (m *usbMIDIDevice) Name() string {
	for _, port := range m.outPorts {
		if port.State == MIDIPortStateConnected {
			return port.Name
		}
	}
	return ""
}

func (m *usbMIDIDevice) GetOutPorts(isConnected bool) ([]MIDIPortInfo, error) {
	res := make([]MIDIPortInfo, 0)
	for _, port := range m.outPorts {
		if port.State == MIDIPortStateConnected || !isConnected {
			res = append(res, port)
		}
	}
	return res, nil
}

// find new connected ports
// TODO: notify linuxsampler for rebind
func (m *usbMIDIDevice) providerEventHandler(devId string, state MIDIPortState) {
	// delete disconnected port
	port, ok := m.outPorts[devId]
	if ok && port.State == MIDIPortStateDisconnected {
		delete(m.outPorts, devId)
		slog.Info("MIDI port disconnected", "devId", devId)
		return
	}
	// new device connected
	if !ok && state == MIDIPortStateConnected {
		ports, err := m.provider.GetMIDIPorts()
		if err != nil {
			slog.Error("failed to get MIDI ports", "error", err)
			return
		}
		for _, p := range ports {
			p.State = state
			m.outPorts[p.DevId] = p
		}
	}
	// TODO: device reconnected: update info
}

// Get MIDI key mapping in format:
// alias : MIDI Key
// eg:
//
//	hihat_close : 42
//	tom1 : 48
func (m *usbMIDIDevice) GetKeysMapping() (map[string]int, error) {
	// TODO: get mapping from repo
	return map[string]int{
		"kick1":            36,
		"snare":            38,
		"snare_rimshot":    39,
		"tom1":             48,
		"tom2":             45,
		"tom3":             43,
		"tom4":             41,
		"hihat_close":      42,
		"hihat_open":       46,
		"hihat_loose":      29,
		"hihat_foot_open":  27,
		"hihat_foot_close": 44,
		"hihat_splash":     28,
		"crash1_edge":      49,
		"ride1_edge":       51,
		"ride1_bell":       53,
	}, nil

}
