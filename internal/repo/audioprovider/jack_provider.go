//go:build linux
// +build linux

package audioprovider

import (
	"fmt"

	"github.com/raspidrum-srv/internal/app/audio"
	"github.com/raspidrum-srv/libs/libalsa"
)

type jackProvider struct {
	device *jackDevice
}

func NewAudioProvider(blackList []string) (audio.AudioDeviceProvider, error) {
	// Get all available cards
	cards, err := libalsa.GetAllCards()
	if err != nil {
		return nil, fmt.Errorf("failed to get ALSA cards: %w", err)
	}

	// Find first suitable card not in blacklist
	for _, cardNum := range cards {
		cardInfo, err := libalsa.GetCardInfo(cardNum)
		if err != nil {
			continue
		}

		// Skip blacklisted devices
		isBlacklisted := false
		for _, blacklisted := range blackList {
			if cardInfo.Name == blacklisted {
				isBlacklisted = true
				break
			}
		}
		if isBlacklisted {
			continue
		}

		// Found suitable card
		return &jackProvider{
			device: &jackDevice{
				cardInfo: cardInfo,
			},
		}, nil
	}

	return nil, fmt.Errorf("no suitable audio devices found")
}

func (p *jackProvider) GetAudioDev() (audio.AudioDevice, error) {
	if p.device == nil {
		return nil, fmt.Errorf("no audio device initialized")
	}
	return p.device, nil
}
