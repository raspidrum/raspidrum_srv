package liblscp

import "strings"

type Effect struct {
	Id          int
	System      string
	Module      string
	Name        string
	Description string
}

func ParseEffect(id int, ln []string) (Effect, error) {
	e := Effect{Id: id}

	for _, v := range ln {
		if vl, f := strings.CutPrefix(v, "SYSTEM: "); f {
			e.System = vl
			continue
		}
		if vl, f := strings.CutPrefix(v, "MODULE: "); f {
			e.Module = vl
			continue
		}
		if vl, f := strings.CutPrefix(v, "NAME: "); f {
			e.Name = vl
			continue
		}
		if vl, f := strings.CutPrefix(v, "DESCRIPTION: "); f {
			e.Description = vl
		}
	}
	return e, nil
}
