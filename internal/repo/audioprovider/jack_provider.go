//go:build linux
// +build linux

package audioprovider

import (
	"context"
	"fmt"
	"time"

	"github.com/raspidrum-srv/internal/app/audio"
	"github.com/raspidrum-srv/internal/repo/dbus"
	"github.com/raspidrum-srv/libs/libalsa"
)

type jackProvider struct {
	device *jackDevice
}

// startJackTransient starts JACK as a transient systemd unit using provided ALSA hw card string (e.g., "hw:0").
func startJackTransient(ctx context.Context, systemd dbus.SystemdManager, hwCard string) error {
	//env := []string{fmt.Sprintf("AUDIO_CARD=%s", hwCard), "JACK_DEFAULT_SERVER=system:playback_1"}
	// exec /usr/bin/jackd -t 2000 -R -P 95 -d alsa -d hw:0,0 -r 48000 -p 512 -n 2 -X seq -s -S
	return systemd.StartTransientUnit(
		ctx,
		"jack.service",
		fmt.Sprintf("JACK Audio Server (%s)", hwCard),
		"/usr/bin/jackd",
		[]string{
			//"-v",
			"-t", "2000",
			"-R",
			"-P", "95",
			"-d", "alsa",
			"-d", hwCard,
			"-r", "48000",
			"-p", "512",
			"-n", "2",
			//"-X", "seq",
			//"-s",
			//"-S",
		},
		[]string{"JACK_NO_AUDIO_RESERVATION=1", "JACK_PROMISCUOUS_SERVER=1"},
		"simple",
	)
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

		// Found suitable card: start JACK transient unit bound to this card
		if sysd, err := dbus.NewDbusSystemdManager(); err == nil {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := startJackTransient(ctx, sysd, fmt.Sprintf("hw:%d,0", cardInfo.ID)); err == nil {
				_ = sysd.WaitForServiceActive(ctx, "jack.service", 5*time.Second)
			} else {
				return nil, err
			}
		}
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
