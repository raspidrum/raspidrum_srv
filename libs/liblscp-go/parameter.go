package liblscp

type ParameterType int

const (
	Unknown ParameterType = iota
	Bool
	Int
	Float
	String
	Bool_list
	Int_list
	Float_list
	String_list
)

var ParameterToName = map[ParameterType]string{
	Bool:   "BOOL",
	Int:    "INT",
	Float:  "FLOAT",
	String: "STRING",
}

var ParameterToType = map[string]ParameterType{
	"BOOL":   Bool,
	"INT":    Int,
	"FLOAT":  Float,
	"STRING": String,
}
