//go:build linux
// +build linux

package audioprovider

import (
	"fmt"

	"github.com/raspidrum-srv/internal/app/audio"
)

type jackProvider struct{}

func NewAudioProvider() (audio.AudioDeviceProvider, error) {
	provider := &jackProvider{}
	return provider, nil
}

func (p *jackProvider) GetAudioDev() (audio.AudioDevice, error) {
	return nil, fmt.Errorf("not implemented")
}
