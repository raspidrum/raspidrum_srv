package liblscp

import (
	"log/slog"
	"strings"
)

type FxSend struct {
	Id             int
	Name           string
	MidiController int
	Level          float64
	AudRouting     []int
	DestChainId    int
	DestChainPos   int
}

func ParseFxSend(multiLineResult []string) (FxSend, error) {
	fs := FxSend{}
	var err error

	for _, v := range multiLineResult {
		if vl, f := strings.CutPrefix(v, "NAME: "); f {
			fs.Name = vl
			continue
		}
		if vl, f := strings.CutPrefix(v, "MIDI_CONTROLLER: "); f {
			fs.MidiController, err = ParseInt(vl)
			if err != nil {
				return fs, err
			}
			continue
		}
		if vl, f := strings.CutPrefix(v, "LEVEL: "); f {
			fs.Level, err = ParseFloat(vl)
			if err != nil {
				return fs, err
			}
			continue
		}
		if vl, f := strings.CutPrefix(v, "AUDIO_OUTPUT_ROUTING: "); f {
			fs.AudRouting, err = ParseIntList(vl)
			if err != nil {
				return fs, err
			}
			continue
		}
		if vl, f := strings.CutPrefix(v, "EFFECT: "); f {
			if vl == "NONE" {
				fs.DestChainId = -1
				fs.DestChainPos = -1
			} else {
				i, err := ParseIntList(vl)
				if err != nil {
					return fs, err
				}
				if len(i) != 2 {
					slog.Info("FxSend: EFFECT field format unknown", slog.String("value", vl))
				} else {
					fs.DestChainId = i[0]
					fs.DestChainPos = i[1]
				}
			}
		}
	}
	return fs, nil
}
