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
	Type           ParameterType
	Value          T
	IsMultiplicity bool
	isMandatory    bool
	Possibilities  []T
	rangeMin       *float64
	rangeMax       *float64
}

func (p *Parameter[T]) SetRange(min *float64, max *float64) {
	p.rangeMin = min
	p.rangeMax = max
}

func (p *Parameter[T]) RangeMin() (setted bool, value float64) {
	if p.rangeMin != nil {
		return true, *p.rangeMin
	}
	return false, 0.0
}

func (p *Parameter[T]) RangeMax() (setted bool, value float64) {
	if p.rangeMax != nil {
		return true, *p.rangeMax
	}
	return false, 0.0
}

func (p *Parameter[T]) GetStringValue() string {
	return fmt.Sprintf("%v", p.Value)
}
