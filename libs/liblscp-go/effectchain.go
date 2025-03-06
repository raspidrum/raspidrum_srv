package liblscp

import "strings"

type EffectChain struct {
	ChainId int
	Effects []EffectInstance
}

func ParseEffectChain(chainId int, lns []string, c *Client) (EffectChain, error) {
	ec := EffectChain{ChainId: chainId}
	for _, v := range lns {
		if vl, f := strings.CutPrefix(v, "EFFECT_SEQUENCE: "); f {
			ids, err := ParseIntList(vl)
			if err != nil {
				return ec, err
			}
			ec.Effects = make([]EffectInstance, len(ids))
			for i, v := range ids {
				ec.Effects[i], err = c.GetEffectInstanceInfo(v)
				if err != nil {
					return ec, err
				}
			}
		}
	}
	return ec, nil
}
