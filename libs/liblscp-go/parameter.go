package liblscp

import (
	"fmt"
)

type ParameterType int

const (
	ptUnknown ParameterType = iota
	ptBool
	ptInt
	ptFloat
	ptString
	ptBool_list
	ptInt_list
	ptFloat_list
	ptString_list
)

var ParameterToName = map[ParameterType]string{
	ptBool:   "BOOL",
	ptInt:    "INT",
	ptFloat:  "FLOAT",
	ptString: "STRING",
}

var ParameterToType = map[string]ParameterType{
	"BOOL":   ptBool,
	"INT":    ptInt,
	"FLOAT":  ptFloat,
	"STRING": ptString,
}

type Parameter[T any] struct {
	Name           string
	Description    string
	Value          T
	IsMultiplicity bool
	isMandatory    bool
	Possibilities  []T
	Default        T
	rangeMin       *float32
	rangeMax       *float32
}

func (p *Parameter[T]) SetRange(min float32, max float32) {
	p.rangeMin = &min
	p.rangeMax = &max
}

func (p *Parameter[T]) SetRangeMin(min float32) {
	p.rangeMin = &min
}

func (p *Parameter[T]) SetRangeMax(max float32) {
	p.rangeMin = &max
}

func (p *Parameter[T]) RangeMin() (setted bool, value float32) {
	if p.rangeMin != nil {
		return true, *p.rangeMin
	}
	return false, 0.0
}

func (p *Parameter[T]) RangeMax() (setted bool, value float32) {
	if p.rangeMax != nil {
		return true, *p.rangeMax
	}
	return false, 0.0
}

func (p *Parameter[T]) GetStringValue() string {
	var res string
	switch any(p.Value).(type) {
	case string:
		res = fmt.Sprintf("'%v'", p.Value)
	default:
		res = fmt.Sprintf("%v", p.Value)
	}
	return res
}
