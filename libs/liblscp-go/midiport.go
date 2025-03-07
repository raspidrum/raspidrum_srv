package liblscp

import (
	"strconv"
	"strings"
)

type MidiPort struct {
	Name   Parameter[string]
	Params []Parameter[string]
}

func ParseMidiPort(multiLineResult []string) (MidiPort, error) {
	mp := MidiPort{}

	for _, v := range multiLineResult {
		if vl, f := strings.CutPrefix(v, "NAME: "); f {
			mp.Name = Parameter[string]{
				Name:  "NAME",
				Value: vl,
			}
			continue
		}
		// any additional parameter
		if bf, af, f := strings.Cut(v, ": "); f {
			vl, _ := strconv.Unquote(af)
			pr := Parameter[string]{
				Name:  bf,
				Value: vl,
			}
			mp.Params = append(mp.Params, pr)
		}
	}
	return mp, nil
}
