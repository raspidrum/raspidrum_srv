//go:build darwin
// +build darwin

package libalsa

import "fmt"

func GetCardInfo(cardNum int) (*AlsaCard, error) {
	return nil, fmt.Errorf("ALSA is not supported on this platform")
}

func GetHardwareInfo(cardNum int) (map[string]string, error) {
	return nil, fmt.Errorf("ALSA is not supported on this platform")
}

func GetAllCards() ([]int, error) {
	return nil, fmt.Errorf("ALSA is not supported on this platform")
}

func ListMidiPorts() ([]MidiPortInfo, error) {
	return nil, fmt.Errorf("ALSA is not supported on this platform")
}
