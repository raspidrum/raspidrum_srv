package liblscp

import "strings"

type EffectInstance struct {
	Id     int
	Params []Parameter[float32]
}

func ParseEffectInstance(id int, ln []string) (EffectInstance, error) {
	ei := EffectInstance{Id: id}
	for _, v := range ln {
		if vl, f := strings.CutPrefix(v, "INPUT_CONTROLS: "); f {
			c, err := parseInt(vl)
			if err != nil {
				return ei, err
			}
			ei.Params = make([]Parameter[float32], c)
		}
	}
	return ei, nil
}

func ParseEffectParameter(ln []string) (Parameter[float32], error) {
	prm := Parameter[float32]{}
	for _, v := range ln {
		if vl, f := strings.CutPrefix(v, "DESCRIPTION: "); f {
			prm.Description = vl
			continue
		}
		if vl, f := strings.CutPrefix(v, "VALUE: "); f {
			f, err := parseFloat(vl)
			if err != nil {
				return prm, err
			}
			prm.Value = f
			continue
		}
		if vl, f := strings.CutPrefix(v, "RANGE_MIN: "); f {
			f, err := parseFloat(vl)
			if err != nil {
				return prm, err
			}
			prm.SetRangeMin(f)
			continue
		}
		if vl, f := strings.CutPrefix(v, "RANGE_MAX: "); f {
			f, err := parseFloat(vl)
			if err != nil {
				return prm, err
			}
			prm.SetRangeMax(f)
			continue
		}
		if vl, f := strings.CutPrefix(v, "DEFAULT: "); f {
			f, err := parseFloat(vl)
			if err != nil {
				return prm, err
			}
			prm.Default = f
			continue
		}
		if vl, f := strings.CutPrefix(v, "POSSIBILITIES: "); f {
			ps, err := parseFloatList(vl)
			if err != nil {
				return prm, err
			}
			prm.Possibilities = ps
		}
	}
	return prm, nil
}
