package liblscp

import (
	"log/slog"
	"strings"
)

type FxSend struct {
	Id             int
	Name           string
	MidiController int
	Level          float32
	AudRouting     []int
	DestChainId    int
	DestChainPos   int
}

func ParseFxSend(id int, multiLineResult []string) (FxSend, error) {
	fs := FxSend{Id: id}
	var err error

	for _, v := range multiLineResult {
		if vl, f := strings.CutPrefix(v, "NAME: "); f {
			fs.Name = vl
			continue
		}
		if vl, f := strings.CutPrefix(v, "MIDI_CONTROLLER: "); f {
			fs.MidiController, err = parseInt(vl)
			if err != nil {
				return fs, err
			}
			continue
		}
		if vl, f := strings.CutPrefix(v, "LEVEL: "); f {
			fs.Level, err = parseFloat(vl)
			if err != nil {
				return fs, err
			}
			continue
		}
		if vl, f := strings.CutPrefix(v, "AUDIO_OUTPUT_ROUTING: "); f {
			fs.AudRouting, err = parseIntList(vl)
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
				i, err := parseIntList(vl)
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
