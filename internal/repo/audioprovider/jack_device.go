package audioprovider

import (
	"fmt"

	"github.com/raspidrum-srv/internal/app/audio"
)

type jackDevice struct{}

func (d *jackDevice) GetChannels(direction audio.AudioChannelDirection) ([]audio.AudioChannel, error) {
	return nil, fmt.Errorf("not implemented")
}

func (d *jackDevice) Driver() string {
	return "JACK"
}

func (d *jackDevice) Name() string {
	return "TODO: implement Name"
}
