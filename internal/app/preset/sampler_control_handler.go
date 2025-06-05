package preset

import (
	"fmt"

	"github.com/raspidrum-srv/internal/repo"
)

type SamplerControlHandler struct {
	sampler         repo.SamplerRepo
	samplerChannels repo.SamplerChannels
}

func NewSamplerControlHandler(sampler repo.SamplerRepo, samplerChannels repo.SamplerChannels) *SamplerControlHandler {
	return &SamplerControlHandler{
		sampler:         sampler,
		samplerChannels: samplerChannels,
	}
}

func (s *SamplerControlHandler) SendChannelMidiCC(channelKey string, cc int, value float32) error {
	chnlId, ok := s.samplerChannels[channelKey]
	if !ok {
		return fmt.Errorf("failed send MIDI CC to channel. invalid channel: %s", channelKey)
	}
	return s.sampler.SendMidiCC(chnlId, cc, value)
}
func (s *SamplerControlHandler) SetChannelVolume(channelKey string, value float32) error {
	chnlId, ok := s.samplerChannels[channelKey]
	if !ok {
		return fmt.Errorf("failed set channel volume. invalid channel: %s", channelKey)
	}
	return s.sampler.SetChannelVolume(chnlId, value)
}
