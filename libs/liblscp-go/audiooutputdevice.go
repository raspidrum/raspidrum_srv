package liblscp

import (
	"strconv"
	"strings"
)

type AudioOutputDevice struct {
	DeviceId   int
	Driver     Parameter[string]
	Channels   Parameter[int]
	Samplerate Parameter[int]
	Active     Parameter[bool]
}

func ParseAudioOutputDevice(deviceId int, multiLineResult []string) (AudioOutputDevice, error) {
	aud := AudioOutputDevice{}

	for _, v := range multiLineResult {
		if vl, f := strings.CutPrefix(v, "CHANNELS: "); f {
			ch, err := parseInt(vl)
			if err != nil {
				return aud, err
			}
			aud.Channels = Parameter[int]{
				Name:  "CHANNELS",
				Value: ch,
			}
			continue
		}
		if vl, f := strings.CutPrefix(v, "SAMPLERATE: "); f {
			sr, err := parseInt(vl)
			if err != nil {
				return aud, err
			}
			aud.Samplerate = Parameter[int]{
				Name:  "SAMPLERATE",
				Value: sr,
			}
			continue
		}
		if vl, f := strings.CutPrefix(v, "ACTIVE: "); f {
			a, err := strconv.ParseBool(vl)
			if err != nil {
				return aud, err
			}
			aud.Active = Parameter[bool]{
				Name:  "ACTIVE",
				Value: a,
			}
			continue
		}
		if vl, f := strings.CutPrefix(v, "DRIVER: "); f {
			aud.Driver = Parameter[string]{
				Name:  "DRIVER",
				Value: vl,
			}
			continue
		}
	}
	return aud, nil
}
