//go:build darwin
// +build darwin

package audioprovider

import (
	"fmt"

	"github.com/raspidrum-srv/internal/app/audio"
)

type coreProvider struct{}

func NewAudioProvider() (audio.AudioDeviceProvider, error) {
	provider := &coreProvider{}
	return provider, nil
}

func (p *coreProvider) GetAudioDev() (audio.AudioDevice, error) {
	return &coreDevice{}, nil
}

type coreDevice struct{}

func (d *coreDevice) GetChannels(direction audio.AudioChannelDirection) ([]audio.AudioChannel, error) {
	return nil, fmt.Errorf("not implemented")
}

func (d *coreDevice) Driver() string {
	return "COREAUDIO"
}

func (d *coreDevice) Name() string {
	return "COREAUDIO"
}
