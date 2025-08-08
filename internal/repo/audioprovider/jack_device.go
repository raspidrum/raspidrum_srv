package audioprovider

import (
	"github.com/raspidrum-srv/internal/app/audio"
	"github.com/raspidrum-srv/libs/libalsa"
)

type jackDevice struct {
	cardInfo *libalsa.AlsaCard
}

func (d *jackDevice) GetChannels(direction audio.AudioChannelDirection) ([]audio.AudioChannel, error) {
	// TODO: implement
	//hwInfo, err := libalsa.GetHardwareInfo(d.cardInfo.ID)
	//if err != nil {
	//	return nil, fmt.Errorf("failed to get hardware info: %w", err)
	//}

	var channels []audio.AudioChannel
	// For now just return empty channels since the interface implementation
	// details need to be worked out with the audio.AudioChannel type
	return channels, nil
}

func (d *jackDevice) Driver() string {
	return "JACK"
}

func (d *jackDevice) Name() string {
	return d.cardInfo.Name
}
